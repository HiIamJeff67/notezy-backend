package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentUserAccountRoutes() {
	userAccountBinder := binders.NewUserAccountBinder()
	userAccountController := controllers.NewUserAccountController(
		services.NewUserAccountService(
			models.NotezyDB,
		),
	)

	userAccountRoutes := DevelopmentRouterGroup.Group("/userAccount")
	userAccountRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
	)
	{
		userAccountRoutes.GET(
			"/getMyAccount",
			userAccountBinder.BindGetMyAccount(
				userAccountController.GetMyAccount,
			),
		)
		userAccountRoutes.PUT(
			"/updateMyAccount",
			userAccountBinder.BindUpdateMyAccount(
				userAccountController.UpdateMyAccount,
			),
		)
	}
}
