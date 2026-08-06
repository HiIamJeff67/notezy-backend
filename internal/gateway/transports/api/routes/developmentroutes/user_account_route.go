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

func configureDevelopmentUserAccountRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userAccountBinder := binders.NewUserAccountBinder()
	userAccountController := controllers.NewUserAccountController(coreClient)

	userAccountRoutes := router.Group("/me/account")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userAccountRoutes.GET(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyAccount"),
					middlewares.ApplyMeterMiddleware("server.requests.userAccount.getMyAccount"),
				},
				defaultMiddlewares,
				userAccountBinder.BindGetMyAccount(userAccountController.GetMyAccount),
			)...,
		)
		userAccountRoutes.PUT(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyAccount"),
					middlewares.ApplyMeterMiddleware("server.requests.userAccount.updateMyAccount"),
				},
				defaultMiddlewares,
				userAccountBinder.BindUpdateMyAccount(userAccountController.UpdateMyAccount),
			)...,
		)
		userAccountRoutes.PUT(
			"/google",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("bindGoogleAccount"),
					middlewares.ApplyMeterMiddleware("server.requests.userAccount.bindGoogleAccount"),
				},
				defaultMiddlewares,
				userAccountBinder.BindBindGoogleAccount(userAccountController.BindGoogleAccount),
			)...,
		)
		userAccountRoutes.DELETE(
			"/google",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("unbindGoogleAccount"),
					middlewares.ApplyMeterMiddleware("server.requests.userAccount.unbindGoogleAccount"),
				},
				defaultMiddlewares,
				userAccountBinder.BindUnbindGoogleAccount(userAccountController.UnbindGoogleAccount),
			)...,
		)
	}
}
