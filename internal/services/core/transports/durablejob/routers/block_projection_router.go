package routers

import (
	"github.com/gin-gonic/gin"

	durablejobdto "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/durablejob/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
)

func ConfigureBlockProjectionRoutes(
	router *gin.Engine,
	blockService services.BlockServiceInterface,
) {
	endpoint := endpoints.NewBlockProjectionEndpoint(blockService)
	router.POST(
		"/durablejob/"+durablejobdto.ApplyBlockProjectionOperation,
		middlewares.DelegationMiddleware(durablejobdto.ApplyBlockProjectionOperation),
		endpoint.Apply,
	)
}
