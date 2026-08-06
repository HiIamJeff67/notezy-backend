package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

func configureDevelopmentUserRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userBinder := binders.NewUserBinder()
	userController := controllers.NewUserController(coreClient)

	userRoutes := router.Group("/users")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userRoutes.GET(
			"/data",
			middlewares.RepositionMiddleware(
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
			middlewares.RepositionMiddleware(
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
			middlewares.RepositionMiddleware(
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
