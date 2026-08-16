package routers

import (
	"github.com/gin-gonic/gin"

	durablejobdto "github.com/HiIamJeff67/notegic-backend/contracts/durable-job/v1"

	blockservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/blocks"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/durablejob/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

func ConfigureBlockProjectionRoutes(
	router *gin.Engine,
	blockService blockservices.BlockServiceInterface,
) {
	endpoint := endpoints.NewBlockProjectionEndpoint(blockService)
	router.POST(
		"/durablejob/"+durablejobdto.ApplyBlockProjectionOperation,
		middlewares.DelegationMiddleware(durablejobdto.ApplyBlockProjectionOperation),
		endpoint.Apply,
	)
}
