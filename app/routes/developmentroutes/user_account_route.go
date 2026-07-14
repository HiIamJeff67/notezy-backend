package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureDevelopmentUserAccountRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userAccountModule := modules.NewUserAccountModule()

	userAccountRoutes := router.Group("/userAccount")
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
		userAccountRoutes.GET(
			"/getMyAccount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyAccount"),
					middlewares.ApplyMeterMiddleware("server.requests.userAccount.getMyAccount"),
				},
				defaultMiddlewares,
				userAccountModule.Binder.BindGetMyAccount(
					userAccountModule.Controller.GetMyAccount,
				),
			)...,
		)
		userAccountRoutes.PUT(
			"/updateMyAccount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyAccount"),
					middlewares.ApplyMeterMiddleware("server.requests.userAccount.updateMyAccount"),
					middlewares.CSRFMiddleware(),
				},
				defaultMiddlewares,
				userAccountModule.Binder.BindUpdateMyAccount(
					userAccountModule.Controller.UpdateMyAccount,
				),
			)...,
		)
		userAccountRoutes.PUT(
			"/bindGoogleAccount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("bindGoogleAccount"),
					middlewares.ApplyMeterMiddleware("server.requests.userAccount.bindGoogleAccount"),
				},
				defaultMiddlewares,
				userAccountModule.Binder.BindBindGoogleAccount(
					userAccountModule.Controller.BindGoogleAccount,
				),
			)...,
		)
		userAccountRoutes.PUT(
			"/unbindGoogleAccount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("unbindGoogleAccount"),
					middlewares.ApplyMeterMiddleware("server.requests.userAccount.unbindGoogleAccount"),
				},
				defaultMiddlewares,
				userAccountModule.Binder.BindUnbindGoogleAccount(
					userAccountModule.Controller.UnbindGoogleAccount,
				),
			)...,
		)
	}
}
