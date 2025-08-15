// app/graphql/server.go
package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	dataloaders "notezy-backend/app/graphql/dataloaders"
	generated "notezy-backend/app/graphql/generated"
	resolvers "notezy-backend/app/graphql/resolvers"
	models "notezy-backend/app/models"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

func GraphQLHandler() gin.HandlerFunc {
	resolver := resolvers.NewResolver(
		dataloaders.NewDataloaders(models.NotezyDB),
		services.NewUserService(
			models.NotezyDB,
		),
		services.NewThemeService(
			models.NotezyDB,
		),
		services.NewShelfService(
			models.NotezyDB,
		),
	)

	config := generated.Config{
		Resolvers: resolver,
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	return func(c *gin.Context) {
		// place the gin.Context into the context.Context
		ctx := context.WithValue(c.Request.Context(), constants.ContextFieldName_Gin_Context, c)
		server.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/graphql")
	return gin.WrapH(h)
}
