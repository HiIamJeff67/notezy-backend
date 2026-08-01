package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func configureDevelopmentBlockRoutes(router *gin.RouterGroup, coreClient *coreadapters.CoreClient) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	blockBinder := binders.NewBlockBinder()
	blockController := controllers.NewBlockController(coreClient)
	blockRoutes := router.Group("/blocks")
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
		"/:blockId",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware("getMyBlockById"),
				middlewares.ApplyMeterMiddleware("server.requests.block.getMyBlockById"),
			},
			append(
				defaultMiddlewares,
				middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
			),
			blockBinder.BindGetMyBlockById(blockController.GetMyBlockById),
		)...,
	)
	blockRoutes.GET(
		"/batch",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware("getMyBlocksByIds"),
				middlewares.ApplyMeterMiddleware("server.requests.block.getMyBlocksByIds"),
			},
			append(
				defaultMiddlewares,
				middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
			),
			blockBinder.BindGetMyBlocksByIds(blockController.GetMyBlocksByIds),
		)...,
	)
	blockRoutes.GET(
		"/block-pack/:blockPackId",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware("getMyBlocksByBlockPackId"),
				middlewares.ApplyMeterMiddleware("server.requests.block.getMyBlocksByBlockPackId"),
			},
			append(
				defaultMiddlewares,
				middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
			),
			blockBinder.BindGetMyBlocksByBlockPackId(blockController.GetMyBlocksByBlockPackId),
		)...,
	)
}
