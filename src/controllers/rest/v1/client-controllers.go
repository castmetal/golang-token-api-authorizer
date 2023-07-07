package rest_controllers_v1

import (
	"net/http"

	"github.com/castmetal/golang-token-api-authorizer/src/config"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/client"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/core/application/dtos"
	use_cases "github.com/castmetal/golang-token-api-authorizer/src/domains/core/use-cases"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/redis"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

// CreateClient  godoc
//
// @Summary  Create a client with CreateClient Data input
// @Description Creating a client with route permissions
// @Tags   Clients
// @Accept   json
// @Produce  json
// @Param   createClient body  dtos.CreateClientDTO true "CreateClient Data"
// @Success  200   {object} dtos.CreateClientResponseDTO
// @Router   /v1/client [post]
func CreateClientControllerV1(c *gin.Context, redisConn redis.IRedisClient, clientRepository client.IClientRepository, createClientRepositories use_cases.ClientRepositories) {
	var createClientDTO dtos.CreateClientDTO
	if err := c.ShouldBindJSON(&createClientDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.Request.Body != nil {
		defer c.Request.Body.Close()
	}

	ucCreateClient, err := use_cases.NewCreateClient(clientRepository, redisConn, createClientRepositories)
	if err != nil {
		errMessage := common.InvalidConnectionError(err.Error())
		common.HandleHttpErrors(errMessage, c)
		return
	}

	response, err := ucCreateClient.Execute(c, &createClientDTO)
	if err != nil {
		common.HandleHttpErrors(err, c)
		return
	}

	c.IndentedJSON(http.StatusCreated, response)
}

// GenerateToken  godoc
//
// @Summary  Generate a JWT token with Generate Token input
// @Description Generates a new JWT token to access and authorize requesting for routes
// @Tags   Token
// @Accept   json
// @Produce  json
// @Param   generateToken body  dtos.GenerateTokenDTO true "Generate Token Data"
// @Success  200   {object} dtos.GenerateTokenResponseDTO
// @Router   /v1/token/generate [post]
func GenerateTokenControllerV1(c *gin.Context, redisConn redis.IRedisClient, clientRepository client.IClientRepository) {
	var generateTokenDTO dtos.GenerateTokenDTO
	if err := c.ShouldBindJSON(&generateTokenDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.Request.Body != nil {
		defer c.Request.Body.Close()
	}

	ucGenerateToken, err := use_cases.NewGenerateToken(clientRepository, redisConn)
	if err != nil {
		errMessage := common.InvalidConnectionError(err.Error())
		common.HandleHttpErrors(errMessage, c)
		return
	}

	response, err := ucGenerateToken.Execute(c, &generateTokenDTO)
	if err != nil {
		common.HandleHttpErrors(err, c)
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

// AlowAcessToken  godoc
//
// @Summary  Concede access to a route
// @Description Allow access to a specific route
// @Tags   Token
// @Accept   json
// @Produce  json
// @Param   allowToken body  dtos.AllowAccessTokenDTO true "Generate Token Data"
// @Success  200   {object} dtos.AllowAccessTokenResponseDTO
// @Router   /v1/token/access [post]
func AllowAccessTokenControllerV1(c *gin.Context, redisConn redis.IRedisClient, clientRepository client.IClientRepository) {
	var allowTokenDTO dtos.AllowAccessTokenDTO
	if err := c.ShouldBindJSON(&allowTokenDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.Request.Body != nil {
		defer c.Request.Body.Close()
	}

	ucAllowAccessToken, err := use_cases.NewAllowAccessToken(clientRepository, redisConn)
	if err != nil {
		errMessage := common.InvalidConnectionError(err.Error())
		common.HandleHttpErrors(errMessage, c)
		return
	}

	response, err := ucAllowAccessToken.Execute(c, &allowTokenDTO)
	if err != nil {
		common.HandleHttpErrors(err, c)
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func SetClientControllers(routerEngine *gin.Engine, config config.EnvStruct, pgConn *pgxpool.Pool, redisConn redis.IRedisClient) {
	clientRepository := repositories.NewClientRepository(pgConn)
	v1 := routerEngine.Group("/v1")
	{
		v1.POST("/client", func(c *gin.Context) {
			createClientRepositories := use_cases.ClientRepositories{
				ScopeRepository:    repositories.NewScopeRepository(pgConn),
				ResourceRepository: repositories.NewResourceRepository(pgConn),
			}
			CreateClientControllerV1(c, redisConn, clientRepository, createClientRepositories)
		})
		v1.POST("/token/generate", func(c *gin.Context) {
			GenerateTokenControllerV1(c, redisConn, clientRepository)
		})
		v1.POST("/token/access", func(c *gin.Context) {
			AllowAccessTokenControllerV1(c, redisConn, clientRepository)
		})
	}

}
