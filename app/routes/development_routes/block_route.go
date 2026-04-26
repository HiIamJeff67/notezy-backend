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

func configureDevelopmentBlockRoutes() {
	blockModule := modules.NewBlockModule()

	blockRoutes := DevelopmentRouterGroup.Group("/block")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		blockRoutes.GET(
			"/getMyBlockById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.GetMyBlockById,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindGetMyBlockById(
					blockModule.Controller.GetMyBlockById,
				),
			)...,
		)
		blockRoutes.GET(
			"/getMyBlocksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlocksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.GetMyBlocksByIds,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindGetMyBlocksByIds(
					blockModule.Controller.GetMyBlocksByIds,
				),
			)...,
		)
		blockRoutes.GET(
			"/getMyBlocksByBlockGroupId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlocksByBlockGroupId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.GetMyBlocksByBlockGroupId,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindGetMyBlocksByBlockGroupId(
					blockModule.Controller.GetMyBlocksByBlockGroupId,
				),
			)...,
		)
		blockRoutes.GET(
			"/getMyBlocksByBlockGroupIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlocksByBlockGroupIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.GetMyBlocksByBlockGroupIds,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindGetMyBlocksByBlockGroupIds(
					blockModule.Controller.GetMyBlocksByBlockGroupIds,
				),
			)...,
		)
		blockRoutes.GET(
			"/getMyBlocksByBlockPackId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlocksByBlockPackId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.GetMyBlocksByBlockPackId,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindGetMyBlocksByBlockPackId(
					blockModule.Controller.GetMyBlocksByBlockPackId,
				),
			)...,
		)
		blockRoutes.GET(
			"/getAllMyBlocks",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyBlocks"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.GetAllMyBlocks,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindGetAllMyBlocks(
					blockModule.Controller.GetAllMyBlocks,
				),
			)...,
		)
		blockRoutes.POST(
			"/insertBlock",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "insertBlock"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.InsertBlock,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindInsertBlock(
					blockModule.Controller.InsertBlock,
				),
			)...,
		)
		blockRoutes.POST(
			"/insertBlocks",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "insertBlocks"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.InsertBlocks,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindInsertBlocks(
					blockModule.Controller.InsertBlocks,
				),
			)...,
		)
		blockRoutes.PUT(
			"/updateMyBlockById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyBlockById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.UpdateMyBlockById,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindUpdateMyBlockById(
					blockModule.Controller.UpdateMyBlockById,
				),
			)...,
		)
		blockRoutes.PUT(
			"/updateMyBlocksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyBlocksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.UpdateMyBlocksByIds,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindUpdateMyBlocksByIds(
					blockModule.Controller.UpdateMyBlocksByIds,
				),
			)...,
		)
		blockRoutes.PATCH(
			"/restoreMyBlockById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyBlockById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.RestoreMyBlockById,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindRestoreMyBlockById(
					blockModule.Controller.RestoreMyBlockById,
				),
			)...,
		)
		blockRoutes.PATCH(
			"/restoreMyBlocksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyBlocksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.RestoreMyBlocksByIds,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindRestoreMyBlocksByIds(
					blockModule.Controller.RestoreMyBlocksByIds,
				),
			)...,
		)
		blockRoutes.DELETE(
			"/deleteMyBlockById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyBlockById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.DeleteMyBlockById,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindDeleteMyBlockById(
					blockModule.Controller.DeleteMyBlockById,
				),
			)...,
		)
		blockRoutes.DELETE(
			"/deleteMyBlocksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyBlocksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Block.DeleteMyBlocksByIds,
					),
				},
				defaultMiddlewares,
				blockModule.Binder.BindDeleteMyBlocksByIds(
					blockModule.Controller.DeleteMyBlocksByIds,
				),
			)...,
		)
	}
}
