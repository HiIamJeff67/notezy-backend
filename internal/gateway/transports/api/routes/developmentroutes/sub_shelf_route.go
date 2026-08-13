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

func configureDevelopmentSubShelfRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreAdapter,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	subShelfBinder := binders.NewSubShelfBinder()
	subShelfController := controllers.NewSubShelfController(coreClient)

	subShelfRoutes := router.Group("/sub-shelves")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		subShelfRoutes.GET(
			"/:sub-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMySubShelfById"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.getMySubShelfById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				subShelfBinder.BindGetMySubShelfById(subShelfController.GetMySubShelfById),
			)...,
		)
		subShelfRoutes.GET(
			"/prev-sub-shelf/:prev-sub-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMySubShelvesByPrevSubShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.getMySubShelvesByPrevSubShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				subShelfBinder.BindGetMySubShelvesByPrevSubShelfId(subShelfController.GetMySubShelvesByPrevSubShelfId),
			)...,
		)
		subShelfRoutes.GET(
			"/root-shelf/:root-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getAllMySubShelvesByRootShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.getAllMySubShelvesByRootShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				subShelfBinder.BindGetAllMySubShelvesByRootShelfId(subShelfController.GetAllMySubShelvesByRootShelfId),
			)...,
		)
		subShelfRoutes.GET(
			"/prev-sub-shelf/:prev-sub-shelf-id/items",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMySubShelvesAndItemsByPrevSubShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.getMySubShelvesAndItemsByPrevSubShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				subShelfBinder.BindGetMySubShelvesAndItemsByPrevSubShelfId(subShelfController.GetMySubShelvesAndItemsByPrevSubShelfId),
			)...,
		)
		subShelfRoutes.POST(
			"/root-shelf/:root-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createSubShelfByRootShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.createSubShelfByRootShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindCreateSubShelfByRootShelfId(subShelfController.CreateSubShelfByRootShelfId),
			)...,
		)
		subShelfRoutes.POST(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createSubShelvesByRootShelfIds"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.CreateSubShelvesByRootShelfIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindCreateSubShelvesByRootShelfIds(subShelfController.CreateSubShelvesByRootShelfIds),
			)...,
		)
		subShelfRoutes.PUT(
			"/:sub-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMySubShelfById"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.updateMySubShelfById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindUpdateMySubShelfById(subShelfController.UpdateMySubShelfById),
			)...,
		)
		subShelfRoutes.PUT(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMySubShelvesByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.UpdateMySubShelvesByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindUpdateMySubShelvesByIds(subShelfController.UpdateMySubShelvesByIds),
			)...,
		)
		subShelfRoutes.PUT(
			"/:sub-shelf-id/position",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("moveMySubShelf"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.moveMySubShelf"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindMoveMySubShelfByRootShelfId(subShelfController.MoveMySubShelfByRootShelfId),
			)...,
		)
		subShelfRoutes.PUT(
			"/position",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("moveMySubShelvesByRootShelfId"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.moveMySubShelvesByRootShelfId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindMoveMySubShelvesByRootShelfId(subShelfController.MoveMySubShelvesByRootShelfId),
			)...,
		)
		subShelfRoutes.PUT(
			"/batch/position",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("moveMySubShelvesByRootShelfIds"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.moveMySubShelvesByRootShelfIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindMoveMySubShelvesByRootShelfIds(subShelfController.MoveMySubShelvesByRootShelfIds),
			)...,
		)
		subShelfRoutes.PATCH(
			"/:sub-shelf-id/restore",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("restoreMySubShelfById"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.restoreMySubShelfById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindRestoreMySubShelfById(subShelfController.RestoreMySubShelfById),
			)...,
		)
		subShelfRoutes.PATCH(
			"/batch/restore",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("restoreMySubShelvesByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.restoreMySubShelvesByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindRestoreMySubShelvesByIds(subShelfController.RestoreMySubShelvesByIds),
			)...,
		)
		subShelfRoutes.DELETE(
			"/:sub-shelf-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMySubShelfById"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.deleteMySubShelfById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindDeleteMySubShelfById(subShelfController.DeleteMySubShelfById),
			)...,
		)
		subShelfRoutes.DELETE(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMySubShelvesByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.subShelf.deleteMySubShelvesByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				subShelfBinder.BindDeleteMySubShelvesByIds(subShelfController.DeleteMySubShelvesByIds),
			)...,
		)
	}
}
