package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"

	binders "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/core/adapters"
)

func configureDevelopmentBlockRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreAdapter,
	rateLimiters RateLimiters,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	blockBinder := binders.NewBlockBinder()
	blockController := controllers.NewBlockController(coreClient)
	blockRoutes := router.Group("/blocks")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(rateLimiters.Unauthorized),
		middlewares.TimeoutMiddleware(3 * time.Second),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		blockRoutes.GET(
			"/:block-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyBlockById"),
					middlewares.ApplyMeterMiddleware("server.requests.block.getMyBlockById"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
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
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				blockBinder.BindGetMyBlocksByIds(blockController.GetMyBlocksByIds),
			)...,
		)
		blockRoutes.GET(
			"/block-pack/:block-pack-id",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyBlocksByBlockPackId"),
					middlewares.ApplyMeterMiddleware("server.requests.block.getMyBlocksByBlockPackId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				blockBinder.BindGetMyBlocksByBlockPackId(blockController.GetMyBlocksByBlockPackId),
			)...,
		)
	}
}
