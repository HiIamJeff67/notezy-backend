package routes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
)

func configureUserInfoRoutes() {
	userInfoRoutes := RouterGroup.Group("/userInfo")
	userInfoRoutes.Use(middlewares.AuthMiddleware())
	{
		userInfoRoutes.GET("/getMyInfo", controllers.GetMyInfo)
	}
}
