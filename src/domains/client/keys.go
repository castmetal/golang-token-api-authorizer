package client

import (
	"crypto/md5"
	"fmt"
)

func GetRedisKeys() map[string]string {

	redisKeys := make(map[string]string, 1)
	redisKeys["REDIS_CLIENT_TOKEN_KEY"] = "/clients/tokens"
	redisKeys["REDIS_CLIENT_KEY"] = "/clients"

	return redisKeys
}

func GetHashKey(apiId string, id string) string {
	hashKey := []byte(fmt.Sprintf("%s-%s", apiId, id))
	hash := md5.Sum(hashKey)
	return fmt.Sprintf("%x", string(hash[:]))
}
