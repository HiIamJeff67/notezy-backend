// app/graphql/server.go
package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/gorm"

	generated "notezy-backend/app/graphql/generated"
	resolvers "notezy-backend/app/graphql/resolvers"
)

func NewGraphQLServer(db *gorm.DB) *handler.Server {
	resolver := resolvers.NewResolver()

	config := generated.Config{
		Resolvers: resolver,
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	return server
}

func NewPlaygroundHandler() http.HandlerFunc {
	return playground.Handler("GraphQL Playground", "/graphql")
}
