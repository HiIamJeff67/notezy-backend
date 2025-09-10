package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentShelfRoutes() {
	shelfBinder := binders.NewShelfBinder()
	shelfController := controllers.NewShelfController(
		services.NewShelfService(
			models.NotezyDB,
		),
	)

	shelfRoutes := DevelopmentRouterGroup.Group("/shelf")
	shelfRoutes.Use(
		middlewares.AuthMiddleware(),
		// middlewares.UserRoleMiddleware(enums.UserRole_Normal),
		middlewares.RateLimitMiddleware(1),
	)
	{
		shelfRoutes.GET(
			"/getMyShelfById",
			shelfBinder.BindGetMyShelfById(
				shelfController.GetMyShelfById,
			),
		)
		shelfRoutes.GET(
			"/searchRecentShelves",
			shelfBinder.BindSearchRecentShelves(
				shelfController.SearchRecentShelves,
			),
		)
		shelfRoutes.POST(
			"/createShelf",
			shelfBinder.BindCreateShelf(
				shelfController.CreateShelf,
			),
		)
		shelfRoutes.PUT(
			"/synchronizeShelves",
			shelfBinder.BindSynchronizeShelves(
				shelfController.SynchronizeShelves,
			),
		)
		shelfRoutes.POST(
			"/restoreMyShelfById",
			shelfBinder.BindRestoreMyShelfById(
				shelfController.RestoreMyShelfById,
			),
		)
		shelfRoutes.POST(
			"/restoreMyShelvesByIds",
			shelfBinder.BindRestoreMyShelvesByIds(
				shelfController.RestoreMyShelvesByIds,
			),
		)
		shelfRoutes.DELETE(
			"/deleteMyShelfById",
			shelfBinder.BindDeleteMyShelfById(
				shelfController.DeleteMyShelfById,
			),
		)
		shelfRoutes.DELETE(
			"/deleteMyShelvesByIds",
			shelfBinder.BindDeleteMyShelvesByIds(
				shelfController.DeleteMyShelvesByIds,
			),
		)
	}
}
