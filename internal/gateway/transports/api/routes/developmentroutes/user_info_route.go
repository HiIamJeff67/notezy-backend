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

func configureDevelopmentUserInfoRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userInfoBinder := binders.NewUserInfoBinder()
	userInfoController := controllers.NewUserInfoController(coreClient)

	userInfoRoutes := router.Group("/me/info")
	defaultsMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userInfoRoutes.GET(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyInfo"),
					middlewares.ApplyMeterMiddleware("server.requests.userInfo.getMyInfo"),
				},
				defaultsMiddlewares,
				userInfoBinder.BindGetMyInfo(userInfoController.GetMyInfo),
			)...,
		)
		userInfoRoutes.PUT(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyInfo"),
					middlewares.ApplyMeterMiddleware("server.requests.userInfo.updateMyInfo"),
				},
				defaultsMiddlewares,
				userInfoBinder.BindUpdateMyInfo(userInfoController.UpdateMyInfo),
			)...,
		)
	}
}
