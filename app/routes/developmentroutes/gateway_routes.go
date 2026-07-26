package developmentroutes

import (
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	"github.com/HiIamJeff67/notezy-backend/app/realtime"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func ConfigureRealtimeRoutes() {
	gateway := realtime.NewGateway()
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
