package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	"notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentUserRoutes() {
	userBinder := binders.NewUserBinder()
	userController := controllers.NewUserController(
		services.NewUserService(
			models.NotezyDB,
		),
	)

	userRoutes := DevelopmentRouterGroup.Group("/user")
	userRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		userRoutes.GET(
			"/getUserData",
			userBinder.BindGetUserData(
				userController.GetUserData,
			),
		)
		userRoutes.GET(
			"/getMe",
			userBinder.BindGetMe(
				userController.GetMe,
			),
		)
		userRoutes.PUT(
			"/updateMe",
			userBinder.BindUpdateMe(
				userController.UpdateMe,
			),
		)
	}
}
