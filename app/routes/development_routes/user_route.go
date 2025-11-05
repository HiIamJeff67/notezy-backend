package developmentroutes

import (
	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentUserRoutes() {
	userModule := modules.NewUserModule()

	userRoutes := DevelopmentRouterGroup.Group("/user")
	userRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		userRoutes.GET(
			"/getUserData",
			userModule.Binder.BindGetUserData(
				userModule.Controller.GetUserData,
			),
		)
		userRoutes.GET(
			"/getMe",
			userModule.Binder.BindGetMe(
				userModule.Controller.GetMe,
			),
		)
		userRoutes.PUT(
			"/updateMe",
			userModule.Binder.BindUpdateMe(
				userModule.Controller.UpdateMe,
			),
		)
	}
}
