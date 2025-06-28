package developmentroutes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
)

func configureDevelopmentUserInfoRoutes() {
	userInfoRoutes := DevelopmentRouterGroup.Group("/userInfo")
	userInfoRoutes.Use(middlewares.AuthMiddleware())
	{
		userInfoRoutes.GET(
			"/getMyInfo",
			controllers.UserInfoController.GetMyInfo,
		)
		userInfoRoutes.PUT(
			"/updateMyInfo",
			controllers.UserInfoController.UpdateMyInfo,
		)
	}
}
