package developmentroutes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentUserRoutes() {
	userController := controllers.NewUserController(
		services.NewUserService(
			models.NotezyDB,
		),
	)

	userRoutes := DevelopmentRouterGroup.Group("/user")
	userRoutes.Use(middlewares.AuthMiddleware())
	{
		userRoutes.GET(
			"/getMe",
			userController.GetMe,
		)
		userRoutes.GET(
			"/all",
			userController.GetAllUsers,
		)
		userRoutes.PATCH(
			"/updateMe",
			userController.UpdateMe,
		)
	}
}
