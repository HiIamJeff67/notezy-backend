package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"

	generated "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/generated"

	resolvers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/graphql/resolvers"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

func GraphQLHandler(coreClient *coreadapters.CoreClient) gin.HandlerFunc {
	resolver := resolvers.NewResolver(coreClient)
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	return func(ctx *gin.Context) {
		requestContext := context.WithValue(
			ctx.Request.Context(),
			sharedcontexts.ContextFieldName_GinContext,
			ctx,
		)
		server.ServeHTTP(ctx.Writer, ctx.Request.WithContext(requestContext))
	}
}

func PlaygroundHandler() gin.HandlerFunc {
	return gin.WrapH(playground.Handler("GraphQL Playground", "/graphql"))
}
