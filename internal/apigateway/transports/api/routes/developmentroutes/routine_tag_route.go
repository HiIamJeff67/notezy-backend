package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	binders "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/core/adapters"
)

func configureDevelopmentRoutineTagRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreAdapter,
	rateLimiters RateLimiters,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	routineTagBinder := binders.NewRoutineTagBinder()
	routineTagController := controllers.NewRoutineTagController(coreClient)

	routineTagRoutes := router.Group("/routine-tags")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(rateLimiters.Unauthorized),
		middlewares.TimeoutMiddleware(3 * time.Second),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		routineTagRoutes.GET(
			"/:routine-tag-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.getMyRoutineTagById"),
				},
				defaultMiddlewares,
				routineTagBinder.BindGetMyRoutineTagById(routineTagController.GetMyRoutineTagById),
			)...,
		)
		routineTagRoutes.GET(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getAllMyRoutineTags"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.getAllMyRoutineTags"),
				},
				defaultMiddlewares,
				routineTagBinder.BindGetAllMyRoutineTags(routineTagController.GetAllMyRoutineTags),
			)...,
		)
		routineTagRoutes.POST(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createRoutineTag"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.createRoutineTag"),
				},
				defaultMiddlewares,
				routineTagBinder.BindCreateRoutineTag(routineTagController.CreateRoutineTag),
			)...,
		)
		routineTagRoutes.POST(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createRoutineTags"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.createRoutineTags"),
				},
				defaultMiddlewares,
				routineTagBinder.BindCreateRoutineTags(routineTagController.CreateRoutineTags),
			)...,
		)
		routineTagRoutes.PUT(
			"/:routine-tag-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.updateMyRoutineTagById"),
				},
				defaultMiddlewares,
				routineTagBinder.BindUpdateMyRoutineTagById(routineTagController.UpdateMyRoutineTagById),
			)...,
		)
		routineTagRoutes.PUT(
			"/batch",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyRoutineTagsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.updateMyRoutineTagsByIds"),
				},
				defaultMiddlewares,
				routineTagBinder.BindUpdateMyRoutineTagsByIds(routineTagController.UpdateMyRoutineTagsByIds),
			)...,
		)
		routineTagRoutes.DELETE(
			"/:routine-tag-id/permanently",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("hardDeleteMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.hardDeleteMyRoutineTagById"),
				},
				defaultMiddlewares,
				routineTagBinder.BindHardDeleteMyRoutineTagById(routineTagController.HardDeleteMyRoutineTagById),
			)...,
		)
		routineTagRoutes.DELETE(
			"/batch/permanently",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("hardDeleteMyRoutineTagsByIds"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.hardDeleteMyRoutineTagsByIds"),
				},
				defaultMiddlewares,
				routineTagBinder.BindHardDeleteMyRoutineTagsByIds(routineTagController.HardDeleteMyRoutineTagsByIds),
			)...,
		)
	}
}
