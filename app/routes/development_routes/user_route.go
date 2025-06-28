package developmentroutes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
)

func configureDevelopmentUserRoutes() {
	userRoutes := DevelopmentRouterGroup.Group("/user")
	userRoutes.Use(middlewares.AuthMiddleware())
	{
		userRoutes.GET(
			"/getMe",
			controllers.UserController.GetMe,
		)
		userRoutes.GET(
			"/all",
			controllers.UserController.GetAllUsers,
		)
		userRoutes.PATCH(
			"/updateMe",
			controllers.UserController.UpdateMe,
		)
	}
}
