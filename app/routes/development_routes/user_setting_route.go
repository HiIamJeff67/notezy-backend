package developmentroutes

import (
	"time"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureUserSettingRoutes() {
	userSettingModule := modules.NewUserSettingModule()

	userSettingRoutes := DevelopmentRouterGroup.Group("/userSetting")
	userSettingRoutes.Use(
		middlewares.TimeoutMiddleware(1*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		userSettingRoutes.GET(
			"/getMySetting",
			userSettingModule.Binder.BindGetMySetting(
				userSettingModule.Controller.GetMySetting,
			),
		)
	}
}
