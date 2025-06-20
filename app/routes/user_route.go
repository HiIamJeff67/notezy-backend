package routes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
)

func configureUserRoutes() {
	userRoutes := RouterGroup.Group("/user")
	userRoutes.Use(middlewares.AuthMiddleware())
	{
		userRoutes.GET("/getMe", controllers.GetMe)
		userRoutes.GET("/all", controllers.GetAllUsers)
		userRoutes.PATCH("/updateMe", controllers.UpdateMe)
	}
}
