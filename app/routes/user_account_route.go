package routes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
)

func configureUserAccountRoutes() {
	userAccountRoutes := RouterGroup.Group("/userAccount")
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
