package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func configureDevelopmentRoutineTagRoutes() {
	routineTagModule := modules.NewRoutineTagModule()

	routineTagRoutes := DevelopmentRouterGroup.Group("/routineTag")
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
		routineTagRoutes.GET(
			"/getMyRoutineTagById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTag.GetMyRoutineTagById,
					),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindGetMyRoutineTagById(
					routineTagModule.Controller.GetMyRoutineTagById,
				),
			)...,
		)
		routineTagRoutes.GET(
			"/getAllMyRoutineTags",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyRoutineTags"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTag.GetAllMyRoutineTags,
					),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindGetAllMyRoutineTags(
					routineTagModule.Controller.GetAllMyRoutineTags,
				),
			)...,
		)
		routineTagRoutes.POST(
			"/createRoutineTag",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createRoutineTag"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTag.CreateRoutineTag,
					),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindCreateRoutineTag(
					routineTagModule.Controller.CreateRoutineTag,
				),
			)...,
		)
		routineTagRoutes.POST(
			"/createRoutineTags",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createRoutineTags"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTag.CreateRoutineTags,
					),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindCreateRoutineTags(
					routineTagModule.Controller.CreateRoutineTags,
				),
			)...,
		)
		routineTagRoutes.PUT(
			"/updateMyRoutineTagById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTag.UpdateMyRoutineTagById,
					),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindUpdateMyRoutineTagById(
					routineTagModule.Controller.UpdateMyRoutineTagById,
				),
			)...,
		)
		routineTagRoutes.PUT(
			"/updateMyRoutineTagsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyRoutineTagsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTag.UpdateMyRoutineTagsByIds,
					),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindUpdateMyRoutineTagsByIds(
					routineTagModule.Controller.UpdateMyRoutineTagsByIds,
				),
			)...,
		)
		routineTagRoutes.DELETE(
			"/hardDeleteMyRoutineTagById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "hardDeleteMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTag.HardDeleteMyRoutineTagById,
					),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindHardDeleteMyRoutineTagById(
					routineTagModule.Controller.HardDeleteMyRoutineTagById,
				),
			)...,
		)
		routineTagRoutes.DELETE(
			"/hardDeleteMyRoutineTagsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "hardDeleteMyRoutineTagsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTag.HardDeleteMyRoutineTagsByIds,
					),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindHardDeleteMyRoutineTagsByIds(
					routineTagModule.Controller.HardDeleteMyRoutineTagsByIds,
				),
			)...,
		)
	}
}
