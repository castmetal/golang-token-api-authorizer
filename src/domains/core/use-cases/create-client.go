package use_cases

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/client"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/core/application/dtos"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/resource"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/scope"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/redis"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/storage/postgres"
	"github.com/google/uuid"
)

type ClientRepositories struct {
	ScopeRepository    scope.IScopeRepository
	ResourceRepository resource.IResourceRepository
}

type (
	CreateClient interface {
		Execute(ctx context.Context, createClientDTO *dtos.CreateClientDTO) (dtos.CreateClientResponseDTO, error)
		toResponse(client *client.Client) dtos.CreateClientResponseDTO
	}
	CreateClientRequest struct {
		CreateClient
		Repository   client.IClientRepository
		Repositories ClientRepositories
		RedisClient  redis.IRedisClient
	}
)

func NewCreateClient(repository client.IClientRepository, redisClient redis.IRedisClient, repositories ClientRepositories) (CreateClient, error) {
	var uc CreateClient = &CreateClientRequest{
		Repository:   repository,
		RedisClient:  redisClient,
		Repositories: repositories,
	}

	return uc, nil
}

func (uc *CreateClientRequest) Execute(ctx context.Context, createClientDTO *dtos.CreateClientDTO) (dtos.CreateClientResponseDTO, error) {
	var response = dtos.CreateClientResponseDTO{}

	_, err := createClientDTO.Validate()
	if err != nil {
		return response, common.InvalidParamsError(err.Error())
	}

	db := uc.Repository.GetDB()
	queries := uc.Repository.GetQueries()
	tx, err := db.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return response, err
	}

	needRollback := true

	defer func() {
		if needRollback {
			tx.Rollback(ctx)
		}
	}()

	qtx := queries.WithTx(tx)

	clientData, err := uc.Repository.FindOneByName(ctx, createClientDTO.ClientName, qtx)
	_, ok := err.(*common.ApplicationError)
	if (err != nil && !ok) || clientData != nil && clientData.ID.String() != "" {
		return response, common.AlreadyExistsError("client name " + createClientDTO.ClientName)
	}

	dtoBytes, err := createClientDTO.ToBytes()
	if err != nil {
		return response, common.DefaultDomainError(err.Error())
	}

	scopeId, err := uc.insertScopeAndResources(ctx, createClientDTO, qtx)
	if err != nil {
		return response, err
	}

	dtoReader := bytes.NewReader(dtoBytes)
	clientProps := getClientProps(dtoReader)

	clientProps.ScopeId = scopeId

	clnt, err := client.NewClientEntity(clientProps)
	if err != nil {
		return response, common.DefaultDomainError(err.Error())
	}

	err = uc.Repository.Create(ctx, clnt, qtx)
	if err != nil {
		return response, err
	}

	tx.Commit(ctx)
	needRollback = false

	return uc.toResponse(clnt), nil
}

func (uc *CreateClientRequest) toResponse(client *client.Client) dtos.CreateClientResponseDTO {
	var response dtos.CreateClientResponseDTO

	clientBytes, err := json.Marshal(client)
	if err != nil {
		return response
	}

	err = json.Unmarshal(clientBytes, &response)
	if err != nil {
		return response
	}

	return response
}

func (uc *CreateClientRequest) insertScopeAndResources(ctx context.Context, createClientDTO *dtos.CreateClientDTO, qtx *postgres.Queries) (uuid.UUID, error) {
	scopeEntity, err := uc.getOrCreateScope(ctx, createClientDTO, qtx)
	if err != nil {
		return uuid.UUID{}, err
	}

	for _, value := range createClientDTO.Permissions {
		_, err := uc.getOrCreateResource(ctx, value, scopeEntity, qtx)
		if err != nil {
			return uuid.UUID{}, err
		}
	}

	return scopeEntity.ID, nil
}

func (uc *CreateClientRequest) getOrCreateScope(ctx context.Context, createClientDTO *dtos.CreateClientDTO, qtx *postgres.Queries) (*scope.Scope, error) {
	scopeData, err := uc.Repositories.ScopeRepository.FindOneByName(ctx, createClientDTO.ScopeName, qtx)
	_, ok := err.(*common.ApplicationError)
	if err != nil && !ok {
		return nil, err
	} else if scopeData != nil && scopeData.ID.String() != "" {
		return scopeData, nil
	}

	var scopeEntity *scope.Scope

	dtoBytes, err := json.Marshal(createClientDTO)
	if err != nil {
		return nil, err
	}

	scopeReader := bytes.NewReader(dtoBytes)
	scopeProps := getScopeProps(scopeReader)
	scopeEntity, err = scope.NewScopeEntity(scopeProps)
	if err != nil {
		return nil, err
	}

	err = uc.Repositories.ScopeRepository.Create(ctx, scopeEntity, qtx)
	if err != nil {
		return nil, err
	}

	return scopeEntity, nil
}

func (uc *CreateClientRequest) getOrCreateResource(ctx context.Context, permission dtos.Permission, scopeEntity *scope.Scope, qtx *postgres.Queries) (*resource.Resource, error) {
	resourceData, err := uc.Repositories.ResourceRepository.FindOneByName(ctx, permission.ResourceName, qtx)
	_, ok := err.(*common.ApplicationError)
	if err != nil && !ok {
		return nil, err
	} else if resourceData != nil && resourceData.ID.String() != "" {
		return resourceData, nil
	}

	resourceData, err = uc.Repositories.ResourceRepository.FindOneByPath(ctx, permission.ResourcePath, permission.ResourceMethod, qtx)
	_, ok = err.(*common.ApplicationError)
	if err != nil && !ok {
		return nil, err
	} else if resourceData != nil && resourceData.ID.String() != "" {
		return resourceData, nil
	}

	b, err := json.Marshal(permission)
	if err != nil {
		return nil, err
	}

	permissionReader := bytes.NewReader(b)
	resourceProps := getResourceProps(permissionReader)

	resourceEntity, err := resource.NewResourceEntity(resourceProps)
	if err != nil {
		return nil, err
	}

	err = uc.Repositories.ResourceRepository.Create(ctx, resourceEntity, qtx)
	if err != nil {
		return nil, err
	}

	_, err = uc.Repositories.ScopeRepository.AttachResourceId(ctx, resourceEntity.ID, scopeEntity.ID, qtx)
	if err != nil {
		return nil, err
	}

	return resourceEntity, nil
}

func getClientProps(message io.Reader) client.ClientProps {
	var clientProps client.ClientProps
	messageBuffer := &bytes.Buffer{}
	messageBuffer.ReadFrom(message)

	if err := json.Unmarshal(messageBuffer.Bytes(), &clientProps); err != nil {
		log.Fatal(err)
	}

	return clientProps
}

func getResourceProps(message io.Reader) resource.ResourceProps {
	var resourceProps resource.ResourceProps
	messageBuffer := &bytes.Buffer{}
	messageBuffer.ReadFrom(message)

	if err := json.Unmarshal(messageBuffer.Bytes(), &resourceProps); err != nil {
		log.Fatal(err)
	}

	return resourceProps
}

func getScopeProps(message io.Reader) scope.ScopeProps {
	var scopeProps scope.ScopeProps
	messageBuffer := &bytes.Buffer{}
	messageBuffer.ReadFrom(message)

	if err := json.Unmarshal(messageBuffer.Bytes(), &scopeProps); err != nil {
		log.Fatal(err)
	}

	return scopeProps
}
