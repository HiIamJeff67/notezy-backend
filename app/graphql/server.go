// app/graphql/server.go
package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
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
	routineScope := scopes.NewRoutineScope()
	rootShelfRepository := repositories.NewRootShelfRepository(scopes.NewRootShelfScope())
	stationRepository := repositories.NewStationRepository(scopes.NewStationScope())
	routineRepository := repositories.NewRoutineRepository(routineScope)
	routineTagRepository := repositories.NewRoutineTagRepository(scopes.NewRoutineTagScope())
	routineTaskRepository := repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope())
	itemRepository := repositories.NewItemRepository(scopes.NewItemScope())
	blockPackRepository := repositories.NewBlockPackRepository(scopes.NewBlockPackScope())
	blockGroupRepository := repositories.NewBlockGroupRepository(scopes.NewBlockGroupScope())
	blockRepository := repositories.NewBlockRepository(scopes.NewBlockScope())
	editableBlockAdapter := adapters.NewEditableBlockAdapter()
	routineTaskPayloadAdapter := adapters.NewRoutineTaskPayloadAdapter(editableBlockAdapter)
	userServices := services.NewUserService(
		models.NotezyDB,
		userRepository,
	)
	themeService := services.NewThemeService(
		models.NotezyDB,
	)
	itemService := services.NewItemService(
		models.NotezyDB,
	)
	blockService := services.NewBlockService(
		models.NotezyDB,
		blockPackRepository,
		blockGroupRepository,
		blockRepository,
		editableBlockAdapter,
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
		routineScope,
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
		routineTaskPayloadAdapter,
	)

	resolver := resolvers.NewResolver(
		dataloaders,
		userServices,
		themeService,
		itemService,
		blockService,
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
