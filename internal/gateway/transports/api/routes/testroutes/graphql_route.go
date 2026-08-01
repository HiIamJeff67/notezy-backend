package testroutes

import (
	graphql "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/graphql"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"

	"github.com/gin-gonic/gin"
)

func ConfigureTestGraphQLRoutes() {
	graphqlRoutes := TestRouterGroup.Group("/graphql")

	graphqlRoutes.Use(
		middlewares.AuthMiddleware(),
		middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
	)
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler(coreadapters.NewConfiguredCoreClient()))
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
