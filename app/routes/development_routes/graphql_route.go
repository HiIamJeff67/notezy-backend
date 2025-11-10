package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	graphql "notezy-backend/app/graphql"
	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
)

func configureDevelopmentGraphQLRoutes() {
	graphqlRoutes := DevelopmentRouterGroup.Group("/graphql")

	graphqlRoutes.Use(
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshAccessTokenInterceptor(),
	)
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler())
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
