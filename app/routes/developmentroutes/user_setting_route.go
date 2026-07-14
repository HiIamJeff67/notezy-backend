package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureUserSettingRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userSettingModule := modules.NewUserSettingModule()

	userSettingRoutes := router.Group("/userSetting")
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
		userSettingRoutes.GET(
			"/getMySetting",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMySetting"),
					middlewares.ApplyMeterMiddleware("server.requests.userSetting.getMySetting"),
				},
				defaultMiddlewares,
				userSettingModule.Binder.BindGetMySetting(
					userSettingModule.Controller.GetMySetting,
				),
			)...,
		)
	}
}
