package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	"notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentUserInfoRoutes() {
	userInfoBinder := binders.NewUserInfoBinder()
	userInfoController := controllers.NewUserInfoController(
		services.NewUserInfoService(
			models.NotezyDB,
		),
	)

	userInfoRoutes := DevelopmentRouterGroup.Group("/userInfo")
	userInfoRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		userInfoRoutes.GET(
			"/getMyInfo",
			userInfoBinder.BindGetMyInfo(
				userInfoController.GetMyInfo,
			),
		)
		userInfoRoutes.PUT(
			"/updateMyInfo",
			userInfoBinder.BindUpdateMyInfo(
				userInfoController.UpdateMyInfo,
			),
		)
	}
}
