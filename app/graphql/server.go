// app/graphql/server.go
package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	caches "github.com/HiIamJeff67/notezy-backend/app/caches"
	dataloaders "github.com/HiIamJeff67/notezy-backend/app/graphql/dataloaders"
	generated "github.com/HiIamJeff67/notezy-backend/app/graphql/generated"
	resolvers "github.com/HiIamJeff67/notezy-backend/app/graphql/resolvers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
	storages "github.com/HiIamJeff67/notezy-backend/app/storages"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

func GraphQLHandler() gin.HandlerFunc {
	dataloaders := dataloaders.NewDataloaders(models.NotezyDB)
	userRepository := repositories.NewUserRepository()
	routineScope := scopes.NewRoutineScope()
	rootShelfScope := scopes.NewRootShelfScope()
	stationScope := scopes.NewStationScope()
	itemScope := scopes.NewItemScope()
	blockScope := scopes.NewBlockScope()
	blockPackScope := scopes.NewBlockPackScope()
	routineTaskScope := scopes.NewRoutineTaskScope()
	subShelfScope := scopes.NewSubShelfScope()
	rootShelfRepository := repositories.NewRootShelfRepository(rootShelfScope)
	usersToShelvesRepository := repositories.NewUsersToShelvesRepository()
	subShelfRepository := repositories.NewSubShelfRepository(subShelfScope)
	stationRepository := repositories.NewStationRepository(stationScope)
	usersToStationsRepository := repositories.NewUsersToStationsRepository()
	routineRepository := repositories.NewRoutineRepository(routineScope)
	routineTagRepository := repositories.NewRoutineTagRepository(scopes.NewRoutineTagScope())
	routineTaskRepository := repositories.NewRoutineTaskRepository(routineTaskScope)
	routineTaskRecordRepository := repositories.NewRoutineTaskRecordRepository(scopes.NewRoutineTaskRecordScope())
	itemRepository := repositories.NewItemRepository(itemScope)
	materialRepository := repositories.NewMaterialRepository(scopes.NewMaterialScope())
	blockPackRepository := repositories.NewBlockPackRepository(blockPackScope)
	blockRepository := repositories.NewBlockRepository(blockScope)
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
		itemScope,
	)
	blockService := services.NewBlockService(
		models.NotezyDB,
		blockScope,
		blockPackScope,
		subShelfScope,
		blockPackRepository,
		blockRepository,
	)
	rootShelfService := services.NewRootShelfService(
		models.NotezyDB,
		rootShelfScope,
		rootShelfRepository,
		usersToShelvesRepository,
		blockPackRepository,
		caches.NewRealtimeLeaseStore(caches.RedisClientMap),
	)
	subShelfService := services.NewSubShelfService(
		models.NotezyDB,
		storages.InMemoryStorage,
		subShelfScope,
		subShelfRepository,
		rootShelfRepository,
		materialRepository,
		blockPackRepository,
		caches.NewRealtimeLeaseStore(caches.RedisClientMap),
	)
	stationService := services.NewStationService(
		models.NotezyDB,
		stationScope,
		stationRepository,
		usersToStationsRepository,
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
		routineTaskScope,
		routineTaskRepository,
		routineTaskPayloadAdapter,
	)
	routineTaskRecordService := services.NewRoutineTaskRecordService(
		models.NotezyDB,
		routineTaskRecordRepository,
	)

	resolver := resolvers.NewResolver(
		dataloaders,
		userServices,
		themeService,
		itemService,
		blockService,
		rootShelfService,
		subShelfService,
		stationService,
		routineService,
		routineTagService,
		routineTaskService,
		routineTaskRecordService,
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
