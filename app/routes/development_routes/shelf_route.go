package developmentroutes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentShelfRoutes() {
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
			shelfController.GetMyShelfById,
		)
		shelfRoutes.GET(
			"/getRecentShelves",
			shelfController.SearchRecentShelves,
		)
		shelfRoutes.POST(
			"/createShelf",
			shelfController.CreateShelf,
		)
		shelfRoutes.PUT(
			"/synchronizeShelves",
			shelfController.SynchronizeShelves,
		)
		shelfRoutes.POST(
			"/restoreMyShelf",
			shelfController.RestoreMyShelfById,
		)
		shelfRoutes.POST(
			"/restoreMyShelves",
			shelfController.RestoreMyShelvesByIds,
		)
		shelfRoutes.DELETE(
			"/deleteMyShelf",
			shelfController.DeleteMyShelfById,
		)
		shelfRoutes.DELETE(
			"/deleteMyShelves",
			shelfController.DeleteMyShelvesByIds,
		)
	}
}
