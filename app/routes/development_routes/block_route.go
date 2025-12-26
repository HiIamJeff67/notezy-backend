package developmentroutes

import (
	"time"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentBlockRoutes() {
	blockModule := modules.NewBlockModule()

	blockRoutes := DevelopmentRouterGroup.Group("/block")
	blockRoutes.Use(
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		blockRoutes.GET(
			"getMyBlockById",
			blockModule.Binder.BindGetMyBlockById(
				blockModule.Controller.GetMyBlockById,
			),
		)
		blockRoutes.GET(
			"getAllMyBlocks",
			blockModule.Binder.BindGetAllMyBlocks(
				blockModule.Controller.GetAllMyBlocks,
			),
		)
	}
}
