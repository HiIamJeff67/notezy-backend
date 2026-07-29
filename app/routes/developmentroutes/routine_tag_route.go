package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureDevelopmentRoutineTagRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	routineTagModule := modules.NewRoutineTagModule()

	routineTagRoutes := router.Group("/routine-tags")
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
			"/:routineTagId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.getMyRoutineTagById"),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindGetMyRoutineTagById(
					routineTagModule.Controller.GetMyRoutineTagById,
				),
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
				routineTagModule.Binder.BindGetAllMyRoutineTags(
					routineTagModule.Controller.GetAllMyRoutineTags,
				),
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
				routineTagModule.Binder.BindCreateRoutineTag(
					routineTagModule.Controller.CreateRoutineTag,
				),
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
				routineTagModule.Binder.BindCreateRoutineTags(
					routineTagModule.Controller.CreateRoutineTags,
				),
			)...,
		)
		routineTagRoutes.PUT(
			"/:routineTagId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.updateMyRoutineTagById"),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindUpdateMyRoutineTagById(
					routineTagModule.Controller.UpdateMyRoutineTagById,
				),
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
				routineTagModule.Binder.BindUpdateMyRoutineTagsByIds(
					routineTagModule.Controller.UpdateMyRoutineTagsByIds,
				),
			)...,
		)
		routineTagRoutes.DELETE(
			"/:routineTagId/permanently",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("hardDeleteMyRoutineTagById"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTag.hardDeleteMyRoutineTagById"),
				},
				defaultMiddlewares,
				routineTagModule.Binder.BindHardDeleteMyRoutineTagById(
					routineTagModule.Controller.HardDeleteMyRoutineTagById,
				),
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
				routineTagModule.Binder.BindHardDeleteMyRoutineTagsByIds(
					routineTagModule.Controller.HardDeleteMyRoutineTagsByIds,
				),
			)...,
		)
	}
}
