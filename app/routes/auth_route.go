package routes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
)

func configureAuthRoutes() {
	authRoutes := RouterGroup.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
		// only protected the logout route
		authRoutes.POST("/logout", middlewares.AuthMiddleware(), controllers.Logout)
	}
}
