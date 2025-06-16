package routes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
)

func configureUserRoutes() {
	userRoutes := RouterGroup.Group("/user")
	userRoutes.Use(middlewares.AuthMiddleware())
	{
		userRoutes.GET("/all", controllers.FindAllUsers)
		userRoutes.PATCH("/updateMe", controllers.UpdateMe)
	}
}
