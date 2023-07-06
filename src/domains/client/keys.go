package client

func GetRedisKeys() map[string]string {

	redisKeys := make(map[string]string, 1)
	redisKeys["REDIS_CLIENT_TOKEN_KEY"] = "/clients/tokens/"

	return redisKeys
}
