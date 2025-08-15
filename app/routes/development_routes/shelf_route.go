package developmentroutes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
	"notezy-backend/app/models"
	"notezy-backend/app/services"
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
