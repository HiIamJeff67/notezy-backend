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
		middlewares.RateLimitMiddleware(1),
	)
	{
		shelfRoutes.POST(
			"/createShelf",
			shelfController.CreateShelf,
		)
		shelfRoutes.PUT(
			"/synchronizeShelves",
			shelfController.SynchronizeShelves,
		)
	}
}
