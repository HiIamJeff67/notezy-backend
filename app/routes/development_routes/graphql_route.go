package developmentroutes

import (
	"github.com/gin-gonic/gin"

	graphql "notezy-backend/app/graphql"
	middlewares "notezy-backend/app/middlewares"
)

func configureDevelopmentGraphQLRoutes() {
	graphqlRoutes := DevelopmentRouterGroup.Group("/graphql")

	graphqlRoutes.Use(middlewares.AuthMiddleware())
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler())
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
