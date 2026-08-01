package developmentroutes

import (
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	websocket "github.com/HiIamJeff67/notezy-backend/internal/gateway/websocket"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
)

func ConfigureRealtimeRoutes() {
	gateway := websocket.NewGateway()
	DevelopmentGatewayRouterGroup := DevelopmentRouter.Group("/" + constants.RealtimeDevelopmentBaseURL)
	DevelopmentGatewayRouterGroup.Use(
		middlewares.DomainWhiteListMiddleware(),
		middlewares.RealtimeUpgradeRateLimitMiddleware(),
	)

	DevelopmentGatewayRouterGroup.GET(
		"",
		gateway.Handle,
	)
}
