package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureDevelopmentBlockRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	blockModule := modules.NewBlockModule()
	blockRoutes := router.Group("/block")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}

	blockRoutes.GET(
		"/getMyBlockById",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware("getMyBlockById"),
				middlewares.ApplyMeterMiddleware("server.requests.block.getMyBlockById"),
			},
			append(
				defaultMiddlewares,
				middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
			),
			blockModule.Binder.BindGetMyBlockById(
				blockModule.Controller.GetMyBlockById,
			),
		)...,
	)
	blockRoutes.GET(
		"/getMyBlocksByIds",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware("getMyBlocksByIds"),
				middlewares.ApplyMeterMiddleware("server.requests.block.getMyBlocksByIds"),
			},
			append(
				defaultMiddlewares,
				middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
			),
			blockModule.Binder.BindGetMyBlocksByIds(
				blockModule.Controller.GetMyBlocksByIds,
			),
		)...,
	)
	blockRoutes.GET(
		"/getMyBlocksByBlockPackId",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware("getMyBlocksByBlockPackId"),
				middlewares.ApplyMeterMiddleware("server.requests.block.getMyBlocksByBlockPackId"),
			},
			append(
				defaultMiddlewares,
				middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
			),
			blockModule.Binder.BindGetMyBlocksByBlockPackId(
				blockModule.Controller.GetMyBlocksByBlockPackId,
			),
		)...,
	)
}
