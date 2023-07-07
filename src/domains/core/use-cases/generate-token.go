package use_cases

import (
	"context"
	"encoding/json"
	"time"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/client"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/core/application/dtos"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/redis"
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

		token, err := client.GenerateTokenJWT(getTokenDuration(clientEntity), clientEntity.ID.String(), []byte("1234"))
		if err != nil {
			return response, err
		}

		response.Token = token

		return response, nil
	}

	return response, nil
}

func getTokenDuration(clientData client.Client) time.Duration {
	switch clientData.KeyPeriod {
	case "days":
		return time.Duration(time.Duration(clientData.KeyTimeDuration) * 24 * time.Hour)
	case "years":
		return time.Duration(time.Duration(clientData.KeyTimeDuration) * 365 * 24 * time.Hour)
	case "minutes":
		return time.Duration(time.Duration(clientData.KeyTimeDuration) * time.Minute)
	case "seconds":
		return time.Duration(time.Duration(clientData.KeyTimeDuration) * time.Second)
	default:
		return time.Duration(time.Duration(clientData.KeyTimeDuration) * 24 * time.Hour)
	}
}
