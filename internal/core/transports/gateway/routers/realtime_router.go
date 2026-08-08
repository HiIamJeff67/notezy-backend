package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/realtime"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureRealtimeRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.RealtimeEndpointInterface,
) {
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
