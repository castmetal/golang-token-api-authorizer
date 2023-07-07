package use_cases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/client"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/core/application/dtos"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/redis"
)

const VALID_TOKEN_MESSAGE = "Token Valid"

type (
	AllowAccessToken interface {
		Execute(ctx context.Context, allowAccessTokenDTO *dtos.AllowAccessTokenDTO) (dtos.AllowAccessTokenResponseDTO, error)
	}
	AllowAccessTokenRequest struct {
		AllowAccessToken
		Repository  client.IClientRepository
		RedisClient redis.IRedisClient
	}
)

func NewAllowAccessToken(repository client.IClientRepository, redisClient redis.IRedisClient) (AllowAccessToken, error) {
	var uc AllowAccessToken = &AllowAccessTokenRequest{
		Repository:  repository,
		RedisClient: redisClient,
	}

	return uc, nil
}

func (uc *AllowAccessTokenRequest) Execute(ctx context.Context, allowAccessTokenDTO *dtos.AllowAccessTokenDTO) (dtos.AllowAccessTokenResponseDTO, error) {
	var response = dtos.AllowAccessTokenResponseDTO{}

	_, err := allowAccessTokenDTO.Validate()
	if err != nil {
		return response, common.InvalidParamsError(err.Error())
	}

	hashStr := client.GetHashKey(allowAccessTokenDTO.ApiId, allowAccessTokenDTO.ClientId)
	key := GetClientKey(hashStr)

	redisResult, _ := uc.RedisClient.GetData(ctx, key)
	if redisResult != "" {
		var clientEntity = client.Client{}
		json.Unmarshal([]byte(redisResult), &clientEntity)

		valid, err := client.ValidateTokenJWT(allowAccessTokenDTO.Token, []byte(clientEntity.Salt))
		if err != nil {
			return response, err
		}

		if !valid {
			return response, fmt.Errorf("Invalid Token")
		}

		response.Message = VALID_TOKEN_MESSAGE

		return response, nil
	}

	return response, nil
}