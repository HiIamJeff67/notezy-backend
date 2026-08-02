package routers

import (
	"github.com/gin-gonic/gin"

	websocketdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/websocket"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/websocket/endpoints"
)

func ConfigureBlockPackRoutes(
	router *gin.RouterGroup,
	realtimeService services.RealtimeServiceInterface,
) {
	endpoint := endpoints.NewBlockPackEndpoint(realtimeService)

	router.POST(
		"/block-pack-channel-permission/validate",
		middlewares.DelegationMiddleware(
			websocketdto.ValidateBlockPackChannelPermissionOperation,
		),
		endpoint.ValidateBlockPackChannelPermission,
	)
}
