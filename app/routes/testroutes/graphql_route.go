package testroutes

import (
	"github.com/HiIamJeff67/notezy-backend/app/graphql"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"

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
