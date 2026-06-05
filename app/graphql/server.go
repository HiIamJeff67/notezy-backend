// app/graphql/server.go
package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	dataloaders "github.com/HiIamJeff67/notezy-backend/app/graphql/dataloaders"
	generated "github.com/HiIamJeff67/notezy-backend/app/graphql/generated"
	resolvers "github.com/HiIamJeff67/notezy-backend/app/graphql/resolvers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

func GraphQLHandler() gin.HandlerFunc {
	dataloaders := dataloaders.NewDataloaders(models.NotezyDB)
	userRepository := repositories.NewUserRepository()
	rootShelfRepository := repositories.NewRootShelfRepository(scopes.NewRootShelfScope())
	stationRepository := repositories.NewStationRepository(scopes.NewStationScope())
	routineRepository := repositories.NewRoutineRepository(scopes.NewRoutineScope())
	routineTagRepository := repositories.NewRoutineTagRepository(scopes.NewRoutineTagScope())
	routineTaskRepository := repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope())
	itemRepository := repositories.NewItemRepository(scopes.NewItemScope())
	userServices := services.NewUserService(
		models.NotezyDB,
		userRepository,
	)
	themeService := services.NewThemeService(
		models.NotezyDB,
	)
	rootShelfService := services.NewRootShelfService(
		models.NotezyDB,
		rootShelfRepository,
	)
	stationService := services.NewStationService(
		models.NotezyDB,
		stationRepository,
	)
	routineService := services.NewRoutineService(
		models.NotezyDB,
		stationRepository,
		routineRepository,
		routineTagRepository,
		routineTaskRepository,
		itemRepository,
	)
	routineTagService := services.NewRoutineTagService(
		models.NotezyDB,
		routineTagRepository,
	)
	routineTaskService := services.NewRoutineTaskService(
		models.NotezyDB,
		routineTaskRepository,
	)

	resolver := resolvers.NewResolver(
		dataloaders,
		userServices,
		themeService,
		rootShelfService,
		stationService,
		routineService,
		routineTagService,
		routineTaskService,
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
		ctx := context.WithValue(c.Request.Context(), types.ContextFieldName_GinContext, c)
		server.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/graphql")
	return gin.WrapH(h)
}
