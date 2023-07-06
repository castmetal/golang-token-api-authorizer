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
// @Summary  Create a client on args input
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
	}

}
