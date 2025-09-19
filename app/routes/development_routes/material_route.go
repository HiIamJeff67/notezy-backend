package developmentroutes

import (
	"notezy-backend/app/adapters"
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
	storages "notezy-backend/app/storages"
)

func configureDevelopmentMaterialRoutes() {
	materialBinder := binders.NewMaterialBinder()
	materialController := controllers.NewMaterialController(
		services.NewMaterialService(
			models.NotezyDB,
			storages.InMemoryStorage,
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
			materialBinder.BindGetMyMaterialById(
				materialController.GetMyMaterialById,
			),
		)
		materialRoutes.POST(
			"/searchMyMaterialsByShelfId",
			materialBinder.BindSearchMyMaterialsByShelfId(
				materialController.SearchMyMaterialsByShelfId,
			),
		)
		materialRoutes.POST(
			"/createTextbookMaterial",
			materialBinder.BindCreateTextbookMaterial(
				materialController.CreateTextbookMaterial,
			),
		)
		materialRoutes.PUT(
			"/saveMyTextbookMaterialById",
			adapters.MultipartAdapter(),
			materialBinder.BindSaveMyTextbookMaterialById(
				materialController.SaveMyTextbookMaterialById,
			),
		)
		materialRoutes.PUT(
			"/moveMyMaterialById",
			materialBinder.BindMoveMyMaterialById(
				materialController.MoveMyMaterialById,
			),
		)
		materialRoutes.PATCH(
			"/restoreMyMaterialById",
			materialBinder.BindRestoreMyMaterialById(
				materialController.RestoreMyMaterialById,
			),
		)
		materialRoutes.PATCH(
			"/restoreMyMaterialsByIds",
			materialBinder.BindRestoreMyMaterialsByIds(
				materialController.RestoreMyMaterialsByIds,
			),
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialById",
			materialBinder.BindDeleteMyMaterialById(
				materialController.DeleteMyMaterialById,
			),
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialsByIds",
			materialBinder.BindDeleteMyMaterialsByIds(
				materialController.DeleteMyMaterialsByIds,
			),
		)
	}
}
