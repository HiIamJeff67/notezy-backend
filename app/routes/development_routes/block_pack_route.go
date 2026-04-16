package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
	metrics "notezy-backend/app/monitor/metrics"
	constants "notezy-backend/shared/constants"
)

func configureDevelopmentBlockPackRoutes() {
	blockPackModule := modules.NewBlockPackModule()

	blockPackRoutes := DevelopmentRouterGroup.Group("/blockPack")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshTokenInterceptor(),
	}
	{
		blockPackRoutes.GET(
			"/getMyBlockPackById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockPackById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.GetMyBlockPackById,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindGetMyBlockPackById(
					blockPackModule.Controller.GetMyBlockPackById,
				),
			)...,
		)
		blockPackRoutes.GET(
			"/getMyBlockPackAndItsParentById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockPackAndItsParentById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.GetMyBlockPackAndItsParentById,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindGetMyBlockPackAndItsParentById(
					blockPackModule.Controller.GetMyBlockPackAndItsParentById,
				),
			)...,
		)
		blockPackRoutes.GET(
			"/getMyBlockPacksByParentSubShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockPacksByParentSubShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.GetMyBlockPacksByParentSubShelfId,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindGetMyBlockPacksByParentSubShelfId(
					blockPackModule.Controller.GetMyBlockPacksByParentSubShelfId,
				),
			)...,
		)
		blockPackRoutes.GET(
			"/getAllMyBlockPacksByRootShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyBlockPacksByRootShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.GetAllMyBlockPacksByRootShelfId,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindGetAllMyBlockPacksByRootShelfId(
					blockPackModule.Controller.GetAllMyBlockPacksByRootShelfId,
				),
			)...,
		)
		blockPackRoutes.POST(
			"/createBlockPack",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createBlockPack"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.CreateBlockPack,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindCreateBlockPack(
					blockPackModule.Controller.CreateBlockPack,
				),
			)...,
		)
		blockPackRoutes.POST(
			"/createBlockPacks",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createBlockPacks"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.CreateBlockPacks,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindCreateBlockPacks(
					blockPackModule.Controller.CreateBlockPacks,
				),
			)...,
		)
		blockPackRoutes.PUT(
			"/updateMyBlockPackById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyBlockPackById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.UpdateMyBlockPackById,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindUpdateMyBlockPackById(
					blockPackModule.Controller.UpdateMyBlockPackById,
				),
			)...,
		)
		blockPackRoutes.PUT(
			"/updateMyBlockPacksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyBlockPacksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.UpdateMyBlockPacksByIds,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindUpdateMyBlockPacksByIds(
					blockPackModule.Controller.UpdateMyBlockPacksByIds,
				),
			)...,
		)
		blockPackRoutes.PUT(
			"/moveMyBlockPackById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMyBlockPackById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.MoveMyBlockPackById,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindMoveMyBlockPackById(
					blockPackModule.Controller.MoveMyBlockPackById,
				),
			)...,
		)
		blockPackRoutes.PUT(
			"/moveMyBlockPacksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMyBlockPacksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.MoveMyBlockPacksByIds,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindMoveMyBlockPacksByIds(
					blockPackModule.Controller.MoveMyBlockPacksByIds,
				),
			)...,
		)
		blockPackRoutes.PUT(
			"/batchMoveMyBlockPacksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "batchMoveMyBlockPacksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.BatchMoveMyBlockPacksByIds,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindBatchMoveMyBlockPacksByIds(
					blockPackModule.Controller.BatchMoveMyBlockPacksByIds,
				),
			)...,
		)
		blockPackRoutes.PATCH(
			"/restoreMyBlockPackById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyBlockPackById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.RestoreMyBlockPackById,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindRestoreMyBlockPackById(
					blockPackModule.Controller.RestoreMyBlockPackById,
				),
			)...,
		)
		blockPackRoutes.PATCH(
			"/restoreMyBlockPacksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyBlockPacksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.RestoreMyBlockPacksByIds,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindRestoreMyBlockPacksByIds(
					blockPackModule.Controller.RestoreMyBlockPacksByIds,
				),
			)...,
		)
		blockPackRoutes.DELETE(
			"/deleteMyBlockPackById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyBlockPackById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.DeleteMyBlockPackById,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindDeleteMyBlockPackById(
					blockPackModule.Controller.DeleteMyBlockPackById,
				),
			)...,
		)
		blockPackRoutes.DELETE(
			"/deleteMyBlockPacksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyBlockPacksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockPack.DeleteMyBlockPacksByIds,
					),
				},
				defaultMiddlewares,
				blockPackModule.Binder.BindDeleteMyBlockPacksByIds(
					blockPackModule.Controller.DeleteMyBlockPacksByIds,
				),
			)...,
		)
	}
}
