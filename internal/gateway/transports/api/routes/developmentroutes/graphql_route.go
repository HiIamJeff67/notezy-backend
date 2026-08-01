package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	graphql "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/graphql"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func configureDevelopmentGraphQLRoutes(router *gin.RouterGroup, coreClient *coreadapters.CoreClient) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	graphqlRoutes := router.Group("/graphql")

	graphqlRoutes.Use(
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	)
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler(coreClient))
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
