package developmentroutes

import (
	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
	"time"
)

func configureDevelopmentUserInfoRoutes() {
	userInfoModule := modules.NewUserInfoModule()

	userInfoRoutes := DevelopmentRouterGroup.Group("/userInfo")
	userInfoRoutes.Use(
		middlewares.TimeoutMiddleware(1*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		userInfoRoutes.GET(
			"/getMyInfo",
			userInfoModule.Binder.BindGetMyInfo(
				userInfoModule.Controller.GetMyInfo,
			),
		)
		userInfoRoutes.PUT(
			"/updateMyInfo",
			userInfoModule.Binder.BindUpdateMyInfo(
				userInfoModule.Controller.UpdateMyInfo,
			),
		)
	}
}
