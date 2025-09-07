package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureUserSettingRoutes() {
	userSettingBinder := binders.NewUserSettingBinder()
	userSettingController := controllers.NewUserSettingController(
		services.NewUserSettingService(
			models.NotezyDB,
		),
	)

	userSettingRoutes := DevelopmentRouterGroup.Group("/userSetting")
	userSettingRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.RateLimitMiddleware(1),
	)
	{
		userSettingRoutes.GET(
			"/getMySetting",
			userSettingBinder.BindGetMySetting(
				userSettingController.GetMySetting,
			),
		)
	}
}
