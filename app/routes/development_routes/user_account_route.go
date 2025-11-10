package developmentroutes

import (
	"time"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentUserAccountRoutes() {
	userAccountModule := modules.NewUserAccountModule()

	userAccountRoutes := DevelopmentRouterGroup.Group("/userAccount")
	userAccountRoutes.Use(
		middlewares.TimeoutMiddleware(1*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		userAccountRoutes.GET(
			"/getMyAccount",
			userAccountModule.Binder.BindGetMyAccount(
				userAccountModule.Controller.GetMyAccount,
			),
		)
		userAccountRoutes.PUT(
			"/updateMyAccount",
			userAccountModule.Binder.BindUpdateMyAccount(
				userAccountModule.Controller.UpdateMyAccount,
			),
		)
	}
}
