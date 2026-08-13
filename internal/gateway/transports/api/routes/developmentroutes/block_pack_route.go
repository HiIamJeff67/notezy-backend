package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"

	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

func configureDevelopmentBlockPackRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreAdapter,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	blockPackBinder := binders.NewBlockPackBinder()
	blockPackController := controllers.NewBlockPackController(coreClient)

	blockPackRoutes := router.Group("/block-packs")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		blockPackRoutes.GET(
			"/:block-pack-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyBlockPackById"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.getMyBlockPackById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				blockPackBinder.BindGetMyBlockPackById(blockPackController.GetMyBlockPackById),
			)...,
		)
		blockPackRoutes.GET(
			"/:block-pack-id/parent",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyBlockPackAndItsParentById"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.getMyBlockPackAndItsParentById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				blockPackBinder.BindGetMyBlockPackAndItsParentById(blockPackController.GetMyBlockPackAndItsParentById),
			)...,
		)
		blockPackRoutes.GET(
			"/sub-shelf/:parent-sub-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyBlockPacksByParentSubShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.getMyBlockPacksByParentSubShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				blockPackBinder.BindGetMyBlockPacksByParentSubShelfId(blockPackController.GetMyBlockPacksByParentSubShelfId),
			)...,
		)
		blockPackRoutes.GET(
			"/root-shelf/:root-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getAllMyBlockPacksByRootShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.getAllMyBlockPacksByRootShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				blockPackBinder.BindGetAllMyBlockPacksByRootShelfId(blockPackController.GetAllMyBlockPacksByRootShelfId),
			)...,
		)
		blockPackRoutes.POST(
			"/sub-shelf/:parent-sub-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createBlockPack"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.createBlockPack"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindCreateBlockPack(blockPackController.CreateBlockPack),
			)...,
		)
		blockPackRoutes.POST(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createBlockPacks"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.CreateBlockPacks"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindCreateBlockPacks(blockPackController.CreateBlockPacks),
			)...,
		)
		blockPackRoutes.PUT(
			"/:block-pack-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyBlockPackById"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.updateMyBlockPackById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindUpdateMyBlockPackById(blockPackController.UpdateMyBlockPackById),
			)...,
		)
		blockPackRoutes.PUT(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyBlockPacksByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.UpdateMyBlockPacksByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindUpdateMyBlockPacksByIds(blockPackController.UpdateMyBlockPacksByIds),
			)...,
		)
		blockPackRoutes.PUT(
			"/:block-pack-id/position",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("moveMyBlockPackById"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.moveMyBlockPackById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindMoveMyBlockPackByParentSubShelfId(blockPackController.MoveMyBlockPackByParentSubShelfId),
			)...,
		)
		blockPackRoutes.PUT(
			"/position",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("moveMyBlockPacksByParentSubShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.moveMyBlockPacksByParentSubShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindMoveMyBlockPacksByParentSubShelfId(blockPackController.MoveMyBlockPacksByParentSubShelfId),
			)...,
		)
		blockPackRoutes.PUT(
			"/batch/position",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("moveMyBlockPacksByParentSubShelfIds"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.moveMyBlockPacksByParentSubShelfIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindMoveMyBlockPacksByParentSubShelfIds(blockPackController.MoveMyBlockPacksByParentSubShelfIds),
			)...,
		)
		blockPackRoutes.PATCH(
			"/:block-pack-id/restore",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("restoreMyBlockPackById"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.restoreMyBlockPackById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindRestoreMyBlockPackById(blockPackController.RestoreMyBlockPackById),
			)...,
		)
		blockPackRoutes.PATCH(
			"/batch/restore",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("restoreMyBlockPacksByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.restoreMyBlockPacksByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindRestoreMyBlockPacksByIds(blockPackController.RestoreMyBlockPacksByIds),
			)...,
		)
		blockPackRoutes.DELETE(
			"/:block-pack-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyBlockPackById"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.deleteMyBlockPackById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindDeleteMyBlockPackById(blockPackController.DeleteMyBlockPackById),
			)...,
		)
		blockPackRoutes.DELETE(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyBlockPacksByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.blockPack.deleteMyBlockPacksByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Write),
				),
				blockPackBinder.BindDeleteMyBlockPacksByIds(blockPackController.DeleteMyBlockPacksByIds),
			)...,
		)
	}
}
