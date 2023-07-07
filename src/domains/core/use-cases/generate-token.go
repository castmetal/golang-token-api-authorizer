package use_cases

import (
	"context"
	"encoding/json"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/client"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/core/application/dtos"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/redis"
	"github.com/google/uuid"
)

type (
	GenerateToken interface {
		Execute(ctx context.Context, generateTokenDTO *dtos.GenerateTokenDTO) (dtos.GenerateTokenResponseDTO, error)
	}
	GenerateTokenRequest struct {
		GenerateToken
		Repository  client.IClientRepository
		RedisClient redis.IRedisClient
	}
)

func NewGenerateToken(repository client.IClientRepository, redisClient redis.IRedisClient) (GenerateToken, error) {
	var uc GenerateToken = &GenerateTokenRequest{
		Repository:  repository,
		RedisClient: redisClient,
	}

	return uc, nil
}

func (uc *GenerateTokenRequest) Execute(ctx context.Context, generateTokenDTO *dtos.GenerateTokenDTO) (dtos.GenerateTokenResponseDTO, error) {
	var response = dtos.GenerateTokenResponseDTO{}

	_, err := generateTokenDTO.Validate()
	if err != nil {
		return response, common.InvalidParamsError(err.Error())
	}

	hashStr := client.GetHashKey(generateTokenDTO.ApiId, generateTokenDTO.ClientId)
	key := GetClientKey(hashStr)

	redisResult, _ := uc.RedisClient.GetData(ctx, key)
	if redisResult != "" {
		var clientEntity = client.Client{}
		json.Unmarshal([]byte(redisResult), &clientEntity)
		token, err := GenerateTokenByClient(&clientEntity)
		if err != nil {
			return response, err
		}

		response.Token = token

		return response, nil
	}

	id, err := uuid.Parse(generateTokenDTO.ClientId)
	if err != nil {
		return response, common.DefaultDomainError("client id ")
	}

	clientData, err := uc.Repository.FindOneById(ctx, id, nil)
	if err != nil || clientData != nil && clientData.ID.String() == "" || clientData == nil {
		return response, common.AlreadyExistsError("client id " + generateTokenDTO.ClientId)
	}

	if clientData.ApiId != generateTokenDTO.ApiId {
		return response, common.InvalidMessageError("access denied")
	}

	token, err := GenerateTokenByClient(clientData)
	if err != nil {
		return response, err
	}

	response.Token = token

	return response, nil
}

func GenerateTokenByClient(clientEntity *client.Client) (string, error) {
	return client.GenerateTokenJWT(clientEntity.KeyPeriod, clientEntity.KeyTimeDuration, clientEntity.ID.String(), []byte(clientEntity.Salt))
}
