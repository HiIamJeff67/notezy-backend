package routes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
)

func configureUserRoutes() {
	userRoutes := RouterGroup.Group("/user")
	userRoutes.Use(middlewares.AuthMiddleware())
	{
		userRoutes.GET("/all", controllers.FindAllUsers)
	}
}
