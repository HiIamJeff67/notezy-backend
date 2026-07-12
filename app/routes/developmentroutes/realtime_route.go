package developmentroutes

import (
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	realtime "github.com/HiIamJeff67/notezy-backend/app/realtime"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func ConfigureRealtimeRoutes() {
	gateway := realtime.NewGateway()

	DevelopmentRouter.GET(
		"/"+constants.RealtimeDevelopmentBaseURL,
		middlewares.DomainWhiteListMiddleware(),
		middlewares.UnauthorizedRateLimitMiddleware(),
		gateway.Handle,
	)
}
