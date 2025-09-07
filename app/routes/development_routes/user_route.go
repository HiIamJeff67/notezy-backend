package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
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
		middlewares.RateLimitMiddleware(1),
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
