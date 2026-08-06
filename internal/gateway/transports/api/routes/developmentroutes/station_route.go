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

func configureDevelopmentStationRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	stationBinder := binders.NewStationBinder()
	stationController := controllers.NewStationController(coreClient)

	stationRoutes := router.Group("/stations")
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
		stationRoutes.GET(
			"/:stationId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyStationById"),
					middlewares.ApplyMeterMiddleware("server.requests.station.getMyStationById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindGetMyStationById(stationController.GetMyStationById),
			)...,
		)
		stationRoutes.GET(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getAllMyStations"),
					middlewares.ApplyMeterMiddleware("server.requests.station.getAllMyStations"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindGetAllMyStations(stationController.GetAllMyStations),
			)...,
		)
		stationRoutes.POST(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createStation"),
					middlewares.ApplyMeterMiddleware("server.requests.station.createStation"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindCreateStation(stationController.CreateStation),
			)...,
		)
		stationRoutes.POST(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createStations"),
					middlewares.ApplyMeterMiddleware("server.requests.station.createStations"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindCreateStations(stationController.CreateStations),
			)...,
		)
		stationRoutes.PUT(
			"/:stationId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyStationById"),
					middlewares.ApplyMeterMiddleware("server.requests.station.updateMyStationById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				stationBinder.BindUpdateMyStationById(stationController.UpdateMyStationById),
			)...,
		)
		stationRoutes.PUT(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyStationsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.station.updateMyStationsByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				stationBinder.BindUpdateMyStationsByIds(stationController.UpdateMyStationsByIds),
			)...,
		)
		stationRoutes.PATCH(
			"/:stationId/restore",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("restoreMyStationById"),
					middlewares.ApplyMeterMiddleware("server.requests.station.restoreMyStationById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Owner),
				),
				stationBinder.BindRestoreMyStationById(stationController.RestoreMyStationById),
			)...,
		)
		stationRoutes.PATCH(
			"/batch/restore",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("restoreMyStationsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.station.restoreMyStationsByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Owner),
				),
				stationBinder.BindRestoreMyStationsByIds(stationController.RestoreMyStationsByIds),
			)...,
		)
		stationRoutes.DELETE(
			"/:stationId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyStationById"),
					middlewares.ApplyMeterMiddleware("server.requests.station.deleteMyStationById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindDeleteMyStationById(stationController.DeleteMyStationById),
			)...,
		)
		stationRoutes.DELETE(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyStationsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.station.deleteMyStationsByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Owner),
				),
				stationBinder.BindDeleteMyStationsByIds(stationController.DeleteMyStationsByIds),
			)...,
		)
		stationRoutes.DELETE(
			"/:stationId/permanently",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("hardDeleteMyStationById"),
					middlewares.ApplyMeterMiddleware("server.requests.station.hardDeleteMyStationById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Owner),
				),
				stationBinder.BindHardDeleteMyStationById(stationController.HardDeleteMyStationById),
			)...,
		)
		stationRoutes.DELETE(
			"/batch/permanently",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("hardDeleteMyStationsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.station.hardDeleteMyStationsByIds"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Owner),
				),
				stationBinder.BindHardDeleteMyStationsByIds(stationController.HardDeleteMyStationsByIds),
			)...,
		)
		/* ============================== Routes for Station Permissions ============================== */

		stationRoutes.GET(
			"/:stationId/permissions/:userPublicId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyStationPermission"),
					middlewares.ApplyMeterMiddleware("server.requests.station.getMyStationPermission"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindGetMyStationPermission(stationController.GetMyStationPermission),
			)...,
		)
		stationRoutes.POST(
			"/:stationId/permissions/:userPublicId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createMyStationPermission"),
					middlewares.ApplyMeterMiddleware("server.requests.station.createMyStationPermission"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				stationBinder.BindCreateMyStationPermission(stationController.CreateMyStationPermission),
			)...,
		)
		stationRoutes.PUT(
			"/:stationId/permissions/:userPublicId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("upsertMyStationPermission"),
					middlewares.ApplyMeterMiddleware("server.requests.station.upsertMyStationPermission"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				stationBinder.BindUpsertMyStationPermission(stationController.UpsertMyStationPermission),
			)...,
		)
		stationRoutes.PUT(
			"/:stationId/permissions",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("upsertMyStationPermissions"),
					middlewares.ApplyMeterMiddleware("server.requests.station.upsertMyStationPermissions"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				stationBinder.BindUpsertMyStationPermissions(stationController.UpsertMyStationPermissions),
			)...,
		)
		stationRoutes.PATCH(
			"/:stationId/permissions/:userPublicId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyStationPermission"),
					middlewares.ApplyMeterMiddleware("server.requests.station.updateMyStationPermission"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				stationBinder.BindUpdateMyStationPermission(stationController.UpdateMyStationPermission),
			)...,
		)
		stationRoutes.POST(
			"/:stationId/ownership",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("transferMyStationOwnership"),
					middlewares.ApplyMeterMiddleware("server.requests.station.transferMyStationOwnership"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Owner),
				),
				stationBinder.BindTransferMyStationOwnership(stationController.TransferMyStationOwnership),
			)...,
		)
		stationRoutes.DELETE(
			"/:stationId/permissions/:userPublicId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyStationPermission"),
					middlewares.ApplyMeterMiddleware("server.requests.station.deleteMyStationPermission"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				stationBinder.BindDeleteMyStationPermission(stationController.DeleteMyStationPermission),
			)...,
		)
		stationRoutes.DELETE(
			"/:stationId/permissions",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyStationPermissions"),
					middlewares.ApplyMeterMiddleware("server.requests.station.deleteMyStationPermissions"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Admin),
				),
				stationBinder.BindDeleteMyStationPermissions(stationController.DeleteMyStationPermissions),
			)...,
		)
		stationRoutes.DELETE(
			"/:stationId/memberships/me",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("leaveMyStation"),
					middlewares.ApplyMeterMiddleware("server.requests.station.leaveMyStation"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindLeaveMyStation(stationController.LeaveMyStation),
			)...,
		)
		stationRoutes.DELETE(
			"/memberships/me",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("leaveMyStations"),
					middlewares.ApplyMeterMiddleware("server.requests.station.leaveMyStations"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindLeaveMyStations(stationController.LeaveMyStations),
			)...,
		)
	}

	/* ============================== Routes for Visualization ============================== */

	visualizationRoutes := router.Group("/stations/visualizations")
	visualizationMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		visualizationRoutes.GET(
			"/total-count",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyTotalCount"),
					middlewares.ApplyMeterMiddleware("server.requests.station.visualizeMyTotalCount"),
				},
				append(
					visualizationMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				stationBinder.BindVisualizeMyTotalCount(stationController.VisualizeMyTotalCount),
			)...,
		)
	}
}
