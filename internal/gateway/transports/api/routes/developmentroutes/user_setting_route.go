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

func configureUserSettingRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreAdapter,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userSettingBinder := binders.NewUserSettingBinder()
	userSettingController := controllers.NewUserSettingController(coreClient)

	userSettingRoutes := router.Group("/me/settings")
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
		userSettingRoutes.GET(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMySetting"),
					middlewares.ApplyMeterMiddleware("server.requests.userSetting.getMySetting"),
				},
				defaultMiddlewares,
				userSettingBinder.BindGetMySetting(userSettingController.GetMySetting),
			)...,
		)
		userSettingRoutes.PUT(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMySetting"),
					middlewares.ApplyMeterMiddleware("server.requests.userSetting.updateMySetting"),
				},
				defaultMiddlewares,
				userSettingBinder.BindUpdateMySetting(userSettingController.UpdateMySetting),
			)...,
		)
	}
}
