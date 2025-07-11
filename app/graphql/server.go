// app/graphql/server.go
package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	generated "notezy-backend/app/graphql/generated"
	resolvers "notezy-backend/app/graphql/resolvers"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
)

func GraphQLHandler() gin.HandlerFunc {
	resolver := resolvers.NewResolver(
		services.NewUserService(
			models.NotezyDB,
		),
	)

	config := generated.Config{
		Resolvers: resolver,
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	return gin.WrapH(server)
}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/graphql")
	return gin.WrapH(h)
}
