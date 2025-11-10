package developmentroutes

import (
	"time"

	adapters "notezy-backend/app/adapters"
	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentMaterialRoutes() {
	materialModule := modules.NewMaterialModule()

	materialRoutes := DevelopmentRouterGroup.Group("/material")
	materialRoutes.Use(
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		// middlewares.UserRoleMiddleware(enums.UserRole_Normal),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		materialRoutes.GET(
			"/getMyMaterialById",
			materialModule.Binder.BindGetMyMaterialById(
				materialModule.Controller.GetMyMaterialById,
			),
		)
		materialRoutes.GET(
			"/getMyMaterialAndItsParentById",
			materialModule.Binder.BindGetMyMaterialAndItsParentById(
				materialModule.Controller.GetMyMaterialAndItsParentById,
			),
		)
		materialRoutes.GET(
			"/getAllMyMaterialsByParentSubShelfId",
			materialModule.Binder.BindGetAllMyMaterialsByParentSubShelfId(
				materialModule.Controller.GetAllMyMaterialsByParentSubShelfId,
			),
		)
		materialRoutes.GET(
			"/getAllMyMaterialsByRootShelfId",
			materialModule.Binder.BindGetAllMyMaterialsByRootShelfId(
				materialModule.Controller.GetAllMyMaterialsByRootShelfId,
			),
		)
		materialRoutes.POST(
			"/createTextbookMaterial",
			materialModule.Binder.BindCreateTextbookMaterial(
				materialModule.Controller.CreateTextbookMaterial,
			),
		)
		materialRoutes.POST(
			"/createNotebookMaterial",
			materialModule.Binder.BindCreateNotebookMaterial(
				materialModule.Controller.CreateNotebookMaterial,
			),
		)
		materialRoutes.PUT(
			"/updateMyMaterialById",
			materialModule.Binder.BindUpdateMyMaterialById(
				materialModule.Controller.UpdateMyMaterialById,
			),
		)
		materialRoutes.PUT(
			"/saveMyNotebookMaterialById",
			adapters.MultipartAdapter(),
			materialModule.Binder.BindSaveMyNotebookMaterialById(
				materialModule.Controller.SaveMyNotebookMaterialById,
			),
		)
		materialRoutes.PUT(
			"/moveMyMaterialById",
			materialModule.Binder.BindMoveMyMaterialById(
				materialModule.Controller.MoveMyMaterialById,
			),
		)
		materialRoutes.PUT(
			"/moveMyMaterialsByIds",
			materialModule.Binder.BindMoveMyMaterialsByIds(
				materialModule.Controller.MoveMyMaterialsByIds,
			),
		)
		materialRoutes.PATCH(
			"/restoreMyMaterialById",
			materialModule.Binder.BindRestoreMyMaterialById(
				materialModule.Controller.RestoreMyMaterialById,
			),
		)
		materialRoutes.PATCH(
			"/restoreMyMaterialsByIds",
			materialModule.Binder.BindRestoreMyMaterialsByIds(
				materialModule.Controller.RestoreMyMaterialsByIds,
			),
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialById",
			materialModule.Binder.BindDeleteMyMaterialById(
				materialModule.Controller.DeleteMyMaterialById,
			),
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialsByIds",
			materialModule.Binder.BindDeleteMyMaterialsByIds(
				materialModule.Controller.DeleteMyMaterialsByIds,
			),
		)
	}
}
