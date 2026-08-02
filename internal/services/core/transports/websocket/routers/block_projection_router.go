package routers

import (
	"github.com/gin-gonic/gin"

	websocketdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/websocket"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/websocket/endpoints"
)

func ConfigureBlockProjectionRoutes(
	router *gin.RouterGroup,
	blockService services.BlockServiceInterface,
) {
	endpoint := endpoints.NewBlockProjectionEndpoint(blockService)

	router.POST(
		"/block-projection/apply",
		middlewares.DelegationMiddleware(
			websocketdto.ApplyBlockProjectionOperation,
		),
		endpoint.ApplyBlockProjection,
	)
}
