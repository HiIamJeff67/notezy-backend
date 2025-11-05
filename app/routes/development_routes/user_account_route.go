package developmentroutes

import (
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentUserAccountRoutes() {
	userAccountModule := modules.NewUserAccountModule()

	userAccountRoutes := DevelopmentRouterGroup.Group("/userAccount")
	userAccountRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
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
