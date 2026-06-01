package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func configureDevelopmentMaterialRoutes() {
	materialModule := modules.NewMaterialModule()

	materialRoutes := DevelopmentRouterGroup.Group("/material")
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
			"/getMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.GetMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindGetMyMaterialById(
					materialModule.Controller.GetMyMaterialById,
				),
			)...,
		)
		materialRoutes.GET(
			"/getMyMaterialAndItsParentById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyMaterialAndItsParentById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.GetMyMaterialAndItsParentById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindGetMyMaterialAndItsParentById(
					materialModule.Controller.GetMyMaterialAndItsParentById,
				),
			)...,
		)
		materialRoutes.GET(
			"/getMyMaterialsByParentSubShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyMaterialsByParentSubShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.GetMyMaterialsByParentSubShelfId,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindGetMyMaterialsByParentSubShelfId(
					materialModule.Controller.GetMyMaterialsByParentSubShelfId,
				),
			)...,
		)
		materialRoutes.GET(
			"/getAllMyMaterialsByRootShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyMaterialsByRootShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.GetAllMyMaterialsByRootShelfId,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindGetAllMyMaterialsByRootShelfId(
					materialModule.Controller.GetAllMyMaterialsByRootShelfId,
				),
			)...,
		)
		materialRoutes.POST(
			"/createMyMaterial",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createMyMaterial"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.CreateMyMaterial,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindCreateMyMaterial(
					materialModule.Controller.CreateMyMaterial,
				),
			)...,
		)
		materialRoutes.PUT(
			"/updateMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.UpdateMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindUpdateMyMaterialById(
					materialModule.Controller.UpdateMyMaterialById,
				),
			)...,
		)
		materialRoutes.PUT(
			"/saveMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "saveMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.SaveMyMaterialById,
					),
					adapters.MultipartAdapter(),
				},
				defaultMiddlewares,
				materialModule.Binder.BindSaveMyMaterialById(
					materialModule.Controller.SaveMyMaterialById,
				),
			)...,
		)
		materialRoutes.PUT(
			"/moveMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.MoveMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindMoveMyMaterialById(
					materialModule.Controller.MoveMyMaterialById,
				),
			)...,
		)
		materialRoutes.PUT(
			"/moveMyMaterialsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.MoveMyMaterialsByIds,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindMoveMyMaterialsByIds(
					materialModule.Controller.MoveMyMaterialsByIds,
				),
			)...,
		)
		materialRoutes.PATCH(
			"/restoreMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.RestoreMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindRestoreMyMaterialById(
					materialModule.Controller.RestoreMyMaterialById,
				),
			)...,
		)
		materialRoutes.PATCH(
			"/restoreMyMaterialsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.RestoreMyMaterialsByIds,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindRestoreMyMaterialsByIds(
					materialModule.Controller.RestoreMyMaterialsByIds,
				),
			)...,
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.DeleteMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindDeleteMyMaterialById(
					materialModule.Controller.DeleteMyMaterialById,
				),
			)...,
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.DeleteMyMaterialsByIds,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindDeleteMyMaterialsByIds(
					materialModule.Controller.DeleteMyMaterialsByIds,
				),
			)...,
		)
	}
}
