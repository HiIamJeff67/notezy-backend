package developmentroutes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func configureDevelopmentMaterialRoutes() {
	materialController := controllers.NewMaterialController(
		services.NewMaterialService(
			models.NotezyDB,
		),
	)

	materialRoutes := DevelopmentRouterGroup.Group("/material")
	materialRoutes.Use(
		middlewares.AuthMiddleware(),
		// middlewares.UserRoleMiddleware(enums.UserRole_Normal),
		middlewares.RateLimitMiddleware(1),
	)
	{
		materialRoutes.GET(
			"/getMyMaterialById",
			materialController.GetMyMaterialById,
		)
		materialRoutes.POST(
			"/searchMyMaterialsByShelfId",
			materialController.SearchMyMaterialsByShelfId,
		)
		materialRoutes.POST(
			"/createTextbookMaterial",
			materialController.CreateTextbookMaterial,
		)
		materialRoutes.PUT(
			"/restoreMyMaterialById",
			materialController.RestoreMyMaterialById,
		)
		materialRoutes.PUT(
			"/restoreMyMaterialsByIds",
			materialController.RestoreMyMaterialsByIds,
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialById",
			materialController.DeleteMyMaterialById,
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialsByIds",
			materialController.DeleteMyMaterialsByIds,
		)
	}
}
