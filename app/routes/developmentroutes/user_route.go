package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureDevelopmentUserRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userModule := modules.NewUserModule()

	userRoutes := router.Group("/user")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userRoutes.GET(
			"/getUserData",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getUserData"),
					middlewares.ApplyMeterMiddleware("server.requests.user.getUserData"),
				},
				defaultMiddlewares,
				userModule.Binder.BindGetUserData(
					userModule.Controller.GetUserData,
				),
			)...,
		)
		userRoutes.GET(
			"/getMe",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMe"),
					middlewares.ApplyMeterMiddleware("server.requests.user.getMe"),
				},
				defaultMiddlewares,
				userModule.Binder.BindGetMe(
					userModule.Controller.GetMe,
				),
			)...,
		)
		userRoutes.PUT(
			"/updateMe",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMe"),
					middlewares.ApplyMeterMiddleware("server.requests.user.updateMe"),
				},
				defaultMiddlewares,
				userModule.Binder.BindUpdateMe(
					userModule.Controller.UpdateMe,
				),
			)...,
		)
	}
}
