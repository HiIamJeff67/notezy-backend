package developmentroutes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentUserAccountRoutes() {
	userAccountController := controllers.NewUserAccountController(
		services.NewUserAccountService(
			models.NotezyDB,
		),
	)

	userAccountRoutes := DevelopmentRouterGroup.Group("/userAccount")
	userAccountRoutes.Use(middlewares.AuthMiddleware())
	{
		userAccountRoutes.GET(
			"/getMyAccount",
			userAccountController.GetMyAccount,
		)
		userAccountRoutes.PUT(
			"/updateMyAccount",
			userAccountController.UpdateMyAccount,
		)
	}
}
