package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureDevelopmentUserInfoRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userInfoModule := modules.NewUserInfoModule()

	userInfoRoutes := router.Group("/userInfo")
	defaultsMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userInfoRoutes.GET(
			"/getMyInfo",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyInfo"),
					middlewares.ApplyMeterMiddleware("server.requests.userInfo.getMyInfo"),
				},
				defaultsMiddlewares,
				userInfoModule.Binder.BindGetMyInfo(
					userInfoModule.Controller.GetMyInfo,
				),
			)...,
		)
		userInfoRoutes.PUT(
			"/updateMyInfo",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("updateMyInfo"),
					middlewares.ApplyMeterMiddleware("server.requests.userInfo.updateMyInfo"),
				},
				defaultsMiddlewares,
				userInfoModule.Binder.BindUpdateMyInfo(
					userInfoModule.Controller.UpdateMyInfo,
				),
			)...,
		)
	}
}
