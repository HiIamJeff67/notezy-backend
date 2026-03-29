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

func configureDevelopmentBlockGroupRoutes() {
	blockGroupModule := modules.NewBlockGroupModule()

	blockGroupRoutes := DevelopmentRouterGroup.Group("/blockGroup")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.TimeoutMiddleware(10 * time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshTokenInterceptor(),
	}
	{
		blockGroupRoutes.GET(
			"/getMyBlockGroupById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockGroupById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.GetMyBlockGroupById,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindGetMyBlockGroupById(
					blockGroupModule.Controller.GetMyBlockGroupById,
				),
			)...,
		)
		blockGroupRoutes.GET(
			"/getMyBlockGroupAndItsBlocksById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockGroupAndItsBlocksById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.GetMyBlockGroupAndItsBlocksById,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindGetMyBlockGroupAndItsBlocksById(
					blockGroupModule.Controller.GetMyBlockGroupAndItsBlocksById,
				),
			)...,
		)
		blockGroupRoutes.GET(
			"/getMyBlockGroupsAndTheirBlocksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockGroupsAndTheirBlocksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.GetMyBlockGroupsAndTheirBlocksByIds,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindGetMyBlockGroupsAndTheirBlocksByIds(
					blockGroupModule.Controller.GetMyBlockGroupsAndTheirBlocksByIds,
				),
			)...,
		)
		blockGroupRoutes.GET(
			"/getMyBlockGroupsAndTheirBlocksByBlockPackId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockGroupsAndTheirBlocksByBlockPackId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.GetMyBlockGroupsAndTheirBlocksByBlockPackId,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindGetMyBlockGroupsAndTheirBlocksByBlockPackId(
					blockGroupModule.Controller.GetMyBlockGroupsAndTheirBlocksByBlockPackId,
				),
			)...,
		)
		blockGroupRoutes.GET(
			"/getMyBlockGroupsByPrevBlockGroupId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockGroupsByPrevBlockGroupId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.GetMyBlockGroupsByPrevBlockGroupId,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindGetMyBlockGroupsByPrevBlockGroupId(
					blockGroupModule.Controller.GetMyBlockGroupsByPrevBlockGroupId,
				),
			)...,
		)
		blockGroupRoutes.GET(
			"/getAllMyBlockGroupsByBlockPackId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyBlockGroupsByBlockPackId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.GetAllMyBlockGroupsByBlockPackId,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindGetAllMyBlockGroupsByBlockPackId(
					blockGroupModule.Controller.GetAllMyBlockGroupsByBlockPackId,
				),
			)...,
		)
		blockGroupRoutes.POST(
			"/insertBlockGroupByBlockPackId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "insertBlockGroupByBlockPackId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.InsertBlockGroupByBlockPackId,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindInsertBlockGroupByBlockPackId(
					blockGroupModule.Controller.InsertBlockGroupByBlockPackId,
				),
			)...,
		)
		blockGroupRoutes.POST(
			"/insertBlockGroupAndItsBlocksByBlockPackId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "insertBlockGroupAndItsBlocksByBlockPackId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.InsertBlockGroupAndItsBlocksByBlockPackId,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindInsertBlockGroupAndItsBlocksByBlockPackId(
					blockGroupModule.Controller.InsertBlockGroupAndItsBlocksByBlockPackId,
				),
			)...,
		)
		blockGroupRoutes.POST(
			"/insertBlockGroupsAndTheirBlocksByBlockPackId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "insertBlockGroupsAndTheirBlocksByBlockPackId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.InsertBlockGroupsAndTheirBlocksByBlockPackId,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindInsertBlockGroupsAndTheirBlocksByBlockPackId(
					blockGroupModule.Controller.InsertBlockGroupsAndTheirBlocksByBlockPackId,
				),
			)...,
		)
		blockGroupRoutes.POST(
			"/insertSequentialBlockGroupsAndTheirBlocksByBlockPackId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "insertSequentialBlockGroupsAndTheirBlocksByBlockPackId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindInsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(
					blockGroupModule.Controller.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId,
				),
			)...,
		)
		blockGroupRoutes.PUT(
			"/moveMyBlockGroupsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMyBlockGroupsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.MoveMyBlockGroupsByIds,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindMoveMyBlockGroupsByIds(
					blockGroupModule.Controller.MoveMyBlockGroupsByIds,
				),
			)...,
		)
		blockGroupRoutes.PATCH(
			"/restoreMyBlockGroupById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyBlockGroupById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.RestoreMyBlockGroupById,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindRestoreMyBlockGroupById(
					blockGroupModule.Controller.RestoreMyBlockGroupById,
				),
			)...,
		)
		blockGroupRoutes.PATCH(
			"/restoreMyBlockGroupsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyBlockGroupsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.RestoreMyBlockGroupsByIds,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindRestoreMyBlockGroupsByIds(
					blockGroupModule.Controller.RestoreMyBlockGroupsByIds,
				),
			)...,
		)
		blockGroupRoutes.DELETE(
			"/deleteMyBlockGroupById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyBlockGroupById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.DeleteMyBlockGroupById,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindDeleteMyBlockGroupById(
					blockGroupModule.Controller.DeleteMyBlockGroupById,
				),
			)...,
		)
		blockGroupRoutes.DELETE(
			"/deleteMyBlockGroupsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyBlockGroupsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.BlockGroup.DeleteMyBlockGroupsByIds,
					),
				},
				defaultMiddlewares,
				blockGroupModule.Binder.BindDeleteMyBlockGroupsByIds(
					blockGroupModule.Controller.DeleteMyBlockGroupsByIds,
				),
			)...,
		)
	}
}
