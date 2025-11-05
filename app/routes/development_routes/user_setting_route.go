package developmentroutes

import (
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureUserSettingRoutes() {
	userSettingModule := modules.NewUserSettingModule()

	userSettingRoutes := DevelopmentRouterGroup.Group("/userSetting")
	userSettingRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
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
