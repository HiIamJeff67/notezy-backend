package developmentroutes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentUserInfoRoutes() {
	userInfoController := controllers.NewUserInfoController(
		services.NewUserInfoService(
			models.NotezyDB,
		),
	)

	userInfoRoutes := DevelopmentRouterGroup.Group("/userInfo")
	userInfoRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.RateLimitMiddleware(1),
	)
	{
		userInfoRoutes.GET(
			"/getMyInfo",
			userInfoController.GetMyInfo,
		)
		userInfoRoutes.PUT(
			"/updateMyInfo",
			userInfoController.UpdateMyInfo,
		)
	}
}
