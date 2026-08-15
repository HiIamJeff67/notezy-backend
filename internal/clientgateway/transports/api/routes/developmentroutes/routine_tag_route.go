package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	binders "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type RoutineTagRouteDependencies struct {
	CoreClient                *coreadapters.CoreAdapter
	AccessTokenCookieHandler  *cookies.CookieHandler
	RefreshTokenCookieHandler *cookies.CookieHandler
	RateLimiters              RateLimiters
}

func configureDevelopmentRoutineTagRoutes(
	router *gin.RouterGroup,
	deps RoutineTagRouteDependencies,
) {
	coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters := deps.CoreClient, deps.AccessTokenCookieHandler, deps.RefreshTokenCookieHandler, deps.RateLimiters
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	routineTagBinder := binders.NewRoutineTagBinder()
	routineTagController := controllers.NewRoutineTagController(coreClient)

	routineTagRoutes := router.Group("/routine-tags")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(rateLimiters.Unauthorized),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.GatewayAuthenticationMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		routineTagRoutes.GET(
			"/:routine-tag-id",
			middlewares.Reposition(
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
			middlewares.Reposition(
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
			middlewares.Reposition(
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
			middlewares.Reposition(
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
			middlewares.Reposition(
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
			middlewares.Reposition(
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
			middlewares.Reposition(
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
			middlewares.Reposition(
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
