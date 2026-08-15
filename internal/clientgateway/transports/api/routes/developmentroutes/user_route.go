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

type UserRouteDependencies struct {
	CoreClient                *coreadapters.CoreAdapter
	AccessTokenCookieHandler  *cookies.CookieHandler
	RefreshTokenCookieHandler *cookies.CookieHandler
	RateLimiters              RateLimiters
}

func configureDevelopmentUserRoutes(
	router *gin.RouterGroup,
	deps UserRouteDependencies,
) {
	coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters := deps.CoreClient, deps.AccessTokenCookieHandler, deps.RefreshTokenCookieHandler, deps.RateLimiters
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userBinder := binders.NewUserBinder()
	userController := controllers.NewUserController(coreClient)

	userRoutes := router.Group("/users")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(rateLimiters.Unauthorized),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.GatewayAuthenticationMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userRoutes.GET(
			"/data",
			middlewares.Reposition(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getUserData"),
					middlewares.ApplyMeterMiddleware("server.requests.user.getUserData"),
				},
				defaultMiddlewares,
				userBinder.BindGetUserData(userController.GetUserData),
			)...,
		)
		userRoutes.GET(
			"/me",
			middlewares.Reposition(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMe"),
					middlewares.ApplyMeterMiddleware("server.requests.user.getMe"),
				},
				defaultMiddlewares,
				userBinder.BindGetMe(userController.GetMe),
			)...,
		)
		userRoutes.PUT(
			"/me",
			middlewares.Reposition(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMe"),
					middlewares.ApplyMeterMiddleware("server.requests.user.updateMe"),
				},
				defaultMiddlewares,
				userBinder.BindUpdateMe(userController.UpdateMe),
			)...,
		)
	}
}
