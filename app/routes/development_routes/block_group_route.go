package developmentroutes

import (
	"time"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentBlockGroupRoutes() {
	blockGroupModule := modules.NewBlockGroupModule()

	blockGroupRoutes := DevelopmentRouterGroup.Group("/blockGroup")
	blockGroupRoutes.Use(
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		blockGroupRoutes.GET(
			"/getMyBlockGroupAndItsBlocksById",
			blockGroupModule.Binder.BindGetMyBlockGroupAndItsBlocksById(
				blockGroupModule.Controller.GetMyBlockGroupAndItsBlocksById,
			),
		)
		blockGroupRoutes.POST(
			"/createBlockGroupAndItsBlocksByBlockPackId",
			blockGroupModule.Binder.BindCreateBlockGroupAndItsBlocksByBlockPackId(
				blockGroupModule.Controller.CreateBlockGroupAndItsBlocksByBlockPackId,
			),
		)
	}
}
