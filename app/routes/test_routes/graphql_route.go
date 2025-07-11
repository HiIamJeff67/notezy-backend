package testroutes

import (
	"notezy-backend/app/graphql"
	middlewares "notezy-backend/app/middlewares"

	"github.com/gin-gonic/gin"
)

func ConfigureTestGraphQLRoutes() {
	graphqlRoutes := TestRouterGroup.Group("/graphql")

	graphqlRoutes.Use(middlewares.AuthMiddleware())
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler())
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
