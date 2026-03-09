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
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshTokenInterceptor(),
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
			middlewares.CSRFMiddleware(),
			userAccountModule.Binder.BindUpdateMyAccount(
				userAccountModule.Controller.UpdateMyAccount,
			),
		)
		userAccountRoutes.PUT(
			"/bindGoogleAccount",
			userAccountModule.Binder.BindBindGoogleAccount(
				userAccountModule.Controller.BindGoogleAccount,
			),
		)
		userAccountRoutes.PUT(
			"/unbindGoogleAccount",
			userAccountModule.Binder.BindUnbindGoogleAccount(
				userAccountModule.Controller.UnbindGoogleAccount,
			),
		)
	}
}
