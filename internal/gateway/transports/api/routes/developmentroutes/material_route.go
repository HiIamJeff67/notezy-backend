package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	adapters "github.com/HiIamJeff67/notezy-backend/internal/adapters"
	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func configureDevelopmentMaterialRoutes(router *gin.RouterGroup, coreClient *coreadapters.CoreClient) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	materialBinder := binders.NewMaterialBinder()
	materialController := controllers.NewMaterialController(coreClient)

	materialRoutes := router.Group("/materials")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		materialRoutes.GET(
			"/:materialId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyMaterialById"),
					middlewares.ApplyMeterMiddleware("server.requests.material.getMyMaterialById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
				),
				materialBinder.BindGetMyMaterialById(materialController.GetMyMaterialById),
			)...,
		)
		materialRoutes.GET(
			"/:materialId/parent",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyMaterialAndItsParentById"),
					middlewares.ApplyMeterMiddleware("server.requests.material.getMyMaterialAndItsParentById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
				),
				materialBinder.BindGetMyMaterialAndItsParentById(materialController.GetMyMaterialAndItsParentById),
			)...,
		)
		materialRoutes.GET(
			"/sub-shelf/:parentSubShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyMaterialsByParentSubShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.material.getMyMaterialsByParentSubShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
				),
				materialBinder.BindGetMyMaterialsByParentSubShelfId(materialController.GetMyMaterialsByParentSubShelfId),
			)...,
		)
		materialRoutes.GET(
			"/root-shelf/:rootShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getAllMyMaterialsByRootShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.material.getAllMyMaterialsByRootShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
				),
				materialBinder.BindGetAllMyMaterialsByRootShelfId(materialController.GetAllMyMaterialsByRootShelfId),
			)...,
		)
		materialRoutes.POST(
			"/sub-shelf/:parentSubShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createMyMaterial"),
					middlewares.ApplyMeterMiddleware("server.requests.material.createMyMaterial"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindCreateMyMaterial(materialController.CreateMyMaterial),
			)...,
		)
		materialRoutes.PUT(
			"/:materialId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyMaterialById"),
					middlewares.ApplyMeterMiddleware("server.requests.material.updateMyMaterialById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindUpdateMyMaterialById(materialController.UpdateMyMaterialById),
			)...,
		)
		materialRoutes.PUT(
			"/:materialId/content",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("saveMyMaterialById"),
					middlewares.ApplyMeterMiddleware("server.requests.material.saveMyMaterialById"),
					adapters.MultipartAdapter(),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindSaveMyMaterialById(materialController.SaveMyMaterialById),
			)...,
		)
		materialRoutes.PUT(
			"/:materialId/parent",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("moveMyMaterialById"),
					middlewares.ApplyMeterMiddleware("server.requests.material.moveMyMaterialById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindMoveMyMaterialById(materialController.MoveMyMaterialById),
			)...,
		)
		materialRoutes.PUT(
			"/batch/parent",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("moveMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.material.moveMyMaterialsByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindMoveMyMaterialsByIds(materialController.MoveMyMaterialsByIds),
			)...,
		)
		materialRoutes.PATCH(
			"/:materialId/restore",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("restoreMyMaterialById"),
					middlewares.ApplyMeterMiddleware("server.requests.material.restoreMyMaterialById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindRestoreMyMaterialById(materialController.RestoreMyMaterialById),
			)...,
		)
		materialRoutes.PATCH(
			"/batch/restore",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("restoreMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.material.restoreMyMaterialsByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindRestoreMyMaterialsByIds(materialController.RestoreMyMaterialsByIds),
			)...,
		)
		materialRoutes.DELETE(
			"/:materialId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyMaterialById"),
					middlewares.ApplyMeterMiddleware("server.requests.material.deleteMyMaterialById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindDeleteMyMaterialById(materialController.DeleteMyMaterialById),
			)...,
		)
		materialRoutes.DELETE(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.material.deleteMyMaterialsByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Write),
				),
				materialBinder.BindDeleteMyMaterialsByIds(materialController.DeleteMyMaterialsByIds),
			)...,
		)
	}
}
