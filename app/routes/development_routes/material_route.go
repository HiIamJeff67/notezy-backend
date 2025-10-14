package developmentroutes

import (
	adapters "notezy-backend/app/adapters"
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	interceptors "notezy-backend/app/interceptors"
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
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		materialRoutes.GET(
			"/getMyMaterialById",
			materialBinder.BindGetMyMaterialById(
				materialController.GetMyMaterialById,
			),
		)
		materialRoutes.GET(
			"/getAllMyMaterialsByParentSubShelfId",
			materialBinder.BindGetAllMyMaterialsByParentSubShelfId(
				materialController.GetAllMyMaterialsByParentSubShelfId,
			),
		)
		materialRoutes.GET(
			"/getAllMyMaterialsByRootShelfId",
			materialBinder.BindGetAllMyMaterialsByRootShelfId(
				materialController.GetAllMyMaterialsByRootShelfId,
			),
		)
		materialRoutes.POST(
			"/createTextbookMaterial",
			materialBinder.BindCreateTextbookMaterial(
				materialController.CreateTextbookMaterial,
			),
		)
		materialRoutes.POST(
			"/createNotebookMaterial",
			materialBinder.BindCreateNotebookMaterial(
				materialController.CreateNotebookMaterial,
			),
		)
		materialRoutes.PUT(
			"/updateMyMaterialById",
			materialBinder.BindUpdateMyMaterialById(
				materialController.UpdateMyMaterialById,
			),
		)
		materialRoutes.PUT(
			"/saveMyNotebookMaterialById",
			adapters.MultipartAdapter(),
			materialBinder.BindSaveMyNotebookMaterialById(
				materialController.SaveMyNotebookMaterialById,
			),
		)
		materialRoutes.PUT(
			"/moveMyMaterialById",
			materialBinder.BindMoveMyMaterialById(
				materialController.MoveMyMaterialById,
			),
		)
		materialRoutes.PUT(
			"/moveMyMaterialsByIds",
			materialBinder.BindMoveMyMaterialsByIds(
				materialController.MoveMyMaterialsByIds,
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
