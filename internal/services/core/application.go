package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	adapters "github.com/HiIamJeff67/notezy-backend/internal/adapters"
	caches "github.com/HiIamJeff67/notezy-backend/internal/caches"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/scopes"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	durablejobrouters "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/durablejob/routers"
	emailtransport "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/email"
	coremiddlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
	gatewayrouters "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/routers"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
	storages "github.com/HiIamJeff67/notezy-backend/internal/shared/storage"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

var ready = make(chan struct{})

func WaitUntilReady() {
	<-ready
}

func NewCoreTransportRouter() *gin.Engine {
	rootShelfScope := scopes.NewRootShelfScope()
	stationScope := scopes.NewStationScope()
	blockScope := scopes.NewBlockScope()
	blockPackScope := scopes.NewBlockPackScope()
	subShelfScope := scopes.NewSubShelfScope()
	materialScope := scopes.NewMaterialScope()
	routineScope := scopes.NewRoutineScope()
	routineTagScope := scopes.NewRoutineTagScope()
	routineTaskScope := scopes.NewRoutineTaskScope()
	routineTaskRecordScope := scopes.NewRoutineTaskRecordScope()
	itemScope := scopes.NewItemScope()

	userRepository := repositories.NewUserRepository()
	userInfoRepository := repositories.NewUserInfoRepository()
	userAccountRepository := repositories.NewUserAccountRepository()
	userSettingRepository := repositories.NewUserSettingRepository()
	rootShelfRepository := repositories.NewRootShelfRepository(rootShelfScope)
	stationRepository := repositories.NewStationRepository(stationScope)
	usersToShelvesRepository := repositories.NewUsersToShelvesRepository()
	usersToStationsRepository := repositories.NewUsersToStationsRepository()
	blockRepository := repositories.NewBlockRepository(blockScope)
	blockPackRepository := repositories.NewBlockPackRepository(blockPackScope)
	subShelfRepository := repositories.NewSubShelfRepository(subShelfScope)
	materialRepository := repositories.NewMaterialRepository(materialScope)
	routineRepository := repositories.NewRoutineRepository(routineScope)
	routineTagRepository := repositories.NewRoutineTagRepository(routineTagScope)
	routineTaskRepository := repositories.NewRoutineTaskRepository(routineTaskScope)
	routineTaskRecordRepository := repositories.NewRoutineTaskRecordRepository(routineTaskRecordScope)
	itemRepository := repositories.NewItemRepository(itemScope)

	oauthService := services.NewOAuthService(config.OAuthGoogleConfig)
	emailClient := emailtransport.NewClient()
	realtimeLeaseStore := caches.NewRealtimeLeaseStore(caches.RedisClientMap)
	editableBlockAdapter := adapters.NewEditableBlockAdapter()
	routineTaskPayloadAdapter := adapters.NewRoutineTaskPayloadAdapter(editableBlockAdapter)

	authService := services.NewAuthService(
		data.NotezyDB,
		userRepository,
		userInfoRepository,
		userAccountRepository,
		userSettingRepository,
		rootShelfRepository,
		oauthService,
		emailClient,
	)
	rootShelfService := services.NewRootShelfService(
		data.NotezyDB,
		rootShelfScope,
		rootShelfRepository,
		usersToShelvesRepository,
		blockPackRepository,
		realtimeLeaseStore,
	)
	stationService := services.NewStationService(
		data.NotezyDB,
		stationScope,
		stationRepository,
		usersToStationsRepository,
	)
	userSettingService := services.NewUserSettingService(
		data.NotezyDB,
		userSettingRepository,
	)
	userInfoService := services.NewUserInfoService(
		data.NotezyDB,
		userInfoRepository,
	)
	userAccountService := services.NewUserAccountService(
		data.NotezyDB,
		userRepository,
		userAccountRepository,
		oauthService,
	)
	userService := services.NewUserService(
		data.NotezyDB,
		userRepository,
	)
	blockService := services.NewBlockService(
		data.NotezyDB,
		blockScope,
		blockPackScope,
		subShelfScope,
		blockPackRepository,
		blockRepository,
	)
	realtimeService := services.NewRealtimeService(
		data.NotezyDB,
		blockPackRepository,
	)
	routineTagService := services.NewRoutineTagService(
		data.NotezyDB,
		routineTagRepository,
	)
	routineTaskRecordService := services.NewRoutineTaskRecordService(
		data.NotezyDB,
		routineTaskRecordRepository,
	)
	subShelfService := services.NewSubShelfService(
		data.NotezyDB,
		storages.InMemoryStorage,
		subShelfScope,
		subShelfRepository,
		rootShelfRepository,
		materialRepository,
		blockPackRepository,
		realtimeLeaseStore,
	)
	blockPackService := services.NewBlockPackService(
		data.NotezyDB,
		blockPackScope,
		subShelfRepository,
		blockPackRepository,
		realtimeLeaseStore,
	)
	materialService := services.NewMaterialService(
		data.NotezyDB,
		storages.InMemoryStorage,
		materialScope,
		subShelfRepository,
		materialRepository,
	)
	routineService := services.NewRoutineService(
		data.NotezyDB,
		routineScope,
		stationRepository,
		routineRepository,
		routineTagRepository,
		routineTaskRepository,
		itemRepository,
	)
	routineTaskService := services.NewRoutineTaskService(
		data.NotezyDB,
		routineTaskScope,
		routineTaskRepository,
		routineTaskPayloadAdapter,
	)
	themeService := services.NewThemeService(data.NotezyDB)
	itemService := services.NewItemService(data.NotezyDB, itemScope)
	badgeService := services.NewBadgeService(data.NotezyDB)
	authMiddleware := coremiddlewares.AuthMiddleware(userRepository)

	router := gatewayrouters.NewRouter(
		authMiddleware,
		authService,
		rootShelfService,
		stationService,
		userSettingService,
		userInfoService,
		userAccountService,
		userService,
		blockService,
		realtimeService,
		routineTagService,
		routineTaskRecordService,
		subShelfService,
		blockPackService,
		materialService,
		routineService,
		routineTaskService,
		themeService,
		itemService,
		badgeService,
	)
	durablejobrouters.ConfigureBlockProjectionRoutes(router, blockService)

	return router
}

func Start() {
	shutdownObservability := observability.Initialize(context.Background())
	defer shutdownObservability()

	data.NotezyDB = data.ConnectToDatabase(config.PostgresDatabaseConfig)
	defer data.DisconnectToDatabase(data.NotezyDB)

	if err := platformredis.DefaultClientManager.ConnectAll(
		caches.UserDataStore.Range,
		caches.RateLimitRecordStore.Range,
		types.Range[int, int]{Start: constants.RealtimeRedisServerNumber, Size: 1},
	); err != nil {
		exception := exceptions.New(
			"ConnectionFailed",
			"Cache",
			"Start",
			"Failed to connect to cache servers",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(
				context.Background(),
				exception.Origin(),
				exception.String(),
			)
		}
		panic(exception)
	}
	defer func() {
		if err := platformredis.DefaultClientManager.DisconnectAll(); err != nil {
			exception := exceptions.New(
				"DisconnectionFailed",
				"Cache",
				"Stop",
				"Failed to disconnect from cache servers",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err)
			if logs.NotezyLogger != nil {
				logs.NotezyLogger.Error(
					context.Background(),
					exception.Origin(),
					exception.String(),
				)
			}
		}
	}()
	reloadRedisLibraries()

	coreTransportListener, err := net.Listen("tcp", config.CoreListenAddress())
	if err != nil {
		fmt.Println("Failed to listen for Core service transport: ", err)
		return
	}
	coreTransportServer := &http.Server{
		Handler: NewCoreTransportRouter(),
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := coreTransportServer.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown Core service transport: ", err)
		}
	}()
	go func() {
		if err := coreTransportServer.Serve(coreTransportListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Failed to serve Core service transport: ", err)
		}
	}()

	close(ready)

	select {}
}

func reloadRedisLibraries() {
	if exception := caches.FlushCacheLibraries(); exception != nil {
		_ = logs.NotezyLogger.JSON(context.Background(), slog.LevelError, exception.String(), exception)
	}
	if exception := caches.LoadRateLimitRecordCacheLibraries(); exception != nil {
		_ = logs.NotezyLogger.JSON(context.Background(), slog.LevelError, exception.String(), exception)
	}
	if exception := caches.LoadUserQuotaCacheLibraries(); exception != nil {
		_ = logs.NotezyLogger.JSON(context.Background(), slog.LevelError, exception.String(), exception)
	}
	// reload other more redis libraries here...
}
