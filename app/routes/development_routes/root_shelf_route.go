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

func configureDevelopmentRootShelfRoutes() {
	rootShelfModule := modules.NewRootShelfModule()

	rootShelfRoutes := DevelopmentRouterGroup.Group("/rootShelf")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshTokenInterceptor(),
	}
	{
		rootShelfRoutes.GET(
			"/getMyRootShelfById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyRootShelfById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.GetMyRootShelfById,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindGetMyRootShelfById(
					rootShelfModule.Controller.GetMyRootShelfById,
				),
			)...,
		)
		rootShelfRoutes.GET(
			"/searchRecentRootShelves",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "searchRecentRootShelves"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.SearchRecentRootShelves,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindSearchRecentRootShelves(
					rootShelfModule.Controller.SearchRecentRootShelves,
				),
			)...,
		)
		rootShelfRoutes.POST(
			"/createRootShelf",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createRootShelf"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.CreateRootShelf,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindCreateRootShelf(
					rootShelfModule.Controller.CreateRootShelf,
				),
			)...,
		)
		rootShelfRoutes.POST(
			"/createRootShelves",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createRootShelves"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.CreateRootShelves,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindCreateRootShelves(
					rootShelfModule.Controller.CreateRootShelves,
				),
			)...,
		)
		rootShelfRoutes.PUT(
			"/updateMyRootShelfById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyRootShelfById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.UpdateMyRootShelfById,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindUpdateMyRootShelfById(
					rootShelfModule.Controller.UpdateMyRootShelfById,
				),
			)...,
		)
		rootShelfRoutes.PUT(
			"/updateMyRootShelvesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyRootShelvesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.UpdateMyRootShelvesByIds,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindUpdateMyRootShelvesByIds(
					rootShelfModule.Controller.UpdateMyRootShelvesByIds,
				),
			)...,
		)
		rootShelfRoutes.PATCH(
			"/restoreMyRootShelfById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyRootShelfById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.RestoreMyRootShelfById,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindRestoreMyRootShelfById(
					rootShelfModule.Controller.RestoreMyRootShelfById,
				),
			)...,
		)
		rootShelfRoutes.PATCH(
			"/restoreMyRootShelvesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyRootShelvesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.RestoreMyRootShelvesByIds,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindRestoreMyRootShelvesByIds(
					rootShelfModule.Controller.RestoreMyRootShelvesByIds,
				),
			)...,
		)
		rootShelfRoutes.DELETE(
			"/deleteMyRootShelfById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyRootShelfById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.DeleteMyRootShelfById,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindDeleteMyRootShelfById(
					rootShelfModule.Controller.DeleteMyRootShelfById,
				),
			)...,
		)
		rootShelfRoutes.DELETE(
			"/deleteMyRootShelvesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyRootShelvesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RootShelf.DeleteMyRootShelvesByIds,
					),
				},
				defaultMiddlewares,
				rootShelfModule.Binder.BindDeleteMyRootShelvesByIds(
					rootShelfModule.Controller.DeleteMyRootShelvesByIds,
				),
			)...,
		)
	}
}
