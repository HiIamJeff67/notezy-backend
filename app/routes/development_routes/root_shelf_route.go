package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentRootShelfRoutes() {
	rootShelfBinder := binders.NewRootShelfBinder()
	rootShelfController := controllers.NewRootShelfController(
		services.NewRootShelfService(
			models.NotezyDB,
		),
	)

	rootShelfRoutes := DevelopmentRouterGroup.Group("/rootShelf")
	rootShelfRoutes.Use(
		middlewares.AuthMiddleware(),
		// middlewares.UserRoleMiddleware(enums.UserRole_Normal),
		middlewares.RateLimitMiddleware(1),
	)
	{
		rootShelfRoutes.GET(
			"/getMyRootShelfById",
			rootShelfBinder.BindGetMyRootShelfById(
				rootShelfController.GetMyRootShelfById,
			),
		)
		rootShelfRoutes.GET(
			"/searchRecentRootShelves",
			rootShelfBinder.BindSearchRecentRootShelves(
				rootShelfController.SearchRecentRootShelves,
			),
		)
		rootShelfRoutes.POST(
			"/createRootShelf",
			rootShelfBinder.BindCreateRootShelf(
				rootShelfController.CreateRootShelf,
			),
		)
		rootShelfRoutes.PUT(
			"/updateMyRootShelfById",
			rootShelfBinder.BindUpdateMyRootShelfById(
				rootShelfController.UpdateMyRootShelfById,
			),
		)
		rootShelfRoutes.DELETE(
			"/deleteMyRootShelfById",
			rootShelfBinder.BindDeleteMyRootShelfById(
				rootShelfController.DeleteMyRootShelfById,
			),
		)
		rootShelfRoutes.DELETE(
			"/deleteMyRootShelvesByIds",
			rootShelfBinder.BindDeleteMyRootShelvesByIds(
				rootShelfController.DeleteMyRootShelvesByIds,
			),
		)
	}
}
