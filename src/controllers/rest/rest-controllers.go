package rest_controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/castmetal/golang-token-api-authorizer/src/config"
	rest_controllers_v1 "github.com/castmetal/golang-token-api-authorizer/src/controllers/rest/v1"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/redis"
)

type RestControllers struct {
	config       config.EnvStruct
	routerEngine *gin.Engine
	pgConn       *pgxpool.Pool
	redisConn    redis.IRedisClient
}

func NewRestControllers(config config.EnvStruct, routerEngine *gin.Engine, pgConn *pgxpool.Pool, redisConn redis.IRedisClient) *RestControllers {
	return &RestControllers{
		config:       config,
		routerEngine: routerEngine,
		pgConn:       pgConn,
		redisConn:    redisConn,
	}
}

func (r *RestControllers) SetRestControllers() {
	switch r.config.ApiVersion {
	case "v1":
		rest_controllers_v1.SetClientControllers(r.routerEngine, r.config, r.pgConn, r.redisConn)
	default:
		rest_controllers_v1.SetClientControllers(r.routerEngine, r.config, r.pgConn, r.redisConn)
	}
}
