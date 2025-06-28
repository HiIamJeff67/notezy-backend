package developmentroutes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
)

func configureDevelopmentUserAccountRoutes() {
	userAccountRoutes := DevelopmentRouterGroup.Group("/userAccount")
	userAccountRoutes.Use(middlewares.AuthMiddleware())
	{
		userAccountRoutes.GET(
			"/getMyAccount",
			controllers.UserAccountController.GetMyAccount,
		)
		userAccountRoutes.PUT(
			"/updateMyAccount",
			controllers.UserAccountController.UpdateMyAccount,
		)
	}
}
