package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/realtime"

	realtimeservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/realtime"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

type RealtimeRouterDependencies struct {
	Service        realtimeservices.RealtimeServiceInterface
	AuthMiddleware gin.HandlerFunc
}

func configureRealtimeRoutes(
	router *gin.RouterGroup,
	deps RealtimeRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	endpoint := endpoints.NewRealtimeEndpoint(deps.Service)
	realtimeRoutes := router.Group("/realtime")
	{
		realtimeRoutes.POST(
			"/connection-ticket/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateMyRealtimeConnectionTicketOperation,
			),
			authMiddleware,
			endpoint.CreateMyRealtimeConnectionTicket,
		)
		realtimeRoutes.POST(
			"/block-pack-channel-ticket/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateMyBlockPackChannelTicketOperation,
			),
			authMiddleware,
			endpoint.CreateMyBlockPackChannelTicket,
		)
	}
}
