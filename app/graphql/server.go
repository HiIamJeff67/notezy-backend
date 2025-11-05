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
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

func GraphQLHandler() gin.HandlerFunc {
	resolver := resolvers.NewResolver(
		dataloaders.NewDataloaders(models.NotezyDB),
		services.NewUserService(
			models.NotezyDB,
			repositories.NewUserRepository(),
		),
		services.NewThemeService(
			models.NotezyDB,
		),
		services.NewRootShelfService(
			models.NotezyDB,
			repositories.NewRootShelfRepository(),
		),
	)

	config := generated.Config{
		Resolvers: resolver,
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	return func(c *gin.Context) {
		// place the gin.Context into the context.Context
		// since we need the fields extracted by the middlewares
		// which are stored in the gin.Context fields,
		// and gin.Context.Request.Context() will not include this part
		ctx := context.WithValue(c.Request.Context(), constants.ContextFieldName_GinContext, c)
		server.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/graphql")
	return gin.WrapH(h)
}
