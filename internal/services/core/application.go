package core

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	userdata "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/cache/userdata"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/scopes"
	storage "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/storage"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	durablejobrouters "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/durablejob/routers"
	emailtransport "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/email"
	coremiddlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
	gatewayrouters "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/routers"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
	workers "github.com/HiIamJeff67/notezy-backend/internal/services/core/workers"
)

func NewCoreTransportRouter() *gin.Engine {
	validator := validation.New()

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
	inMemoryStorage := storage.NewInMemoryStorage()

	oauthService := services.NewOAuthService(config.OAuthGoogleConfig)
	emailClient := emailtransport.NewClient()

	authService := services.NewAuthService(
		validator,
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
		validator,
		data.NotezyDB,
		rootShelfScope,
		rootShelfRepository,
		usersToShelvesRepository,
		blockPackRepository,
	)
	stationService := services.NewStationService(
		validator,
		data.NotezyDB,
		stationScope,
		stationRepository,
		usersToStationsRepository,
	)
	userSettingService := services.NewUserSettingService(
		validator,
		data.NotezyDB,
		userSettingRepository,
	)
	userInfoService := services.NewUserInfoService(
		validator,
		data.NotezyDB,
		userInfoRepository,
	)
	userAccountService := services.NewUserAccountService(
		validator,
		data.NotezyDB,
		userRepository,
		userAccountRepository,
		oauthService,
	)
	userService := services.NewUserService(
		validator,
		data.NotezyDB,
		userRepository,
	)
	blockService := services.NewBlockService(
		validator,
		data.NotezyDB,
		blockScope,
		blockPackScope,
		subShelfScope,
		blockPackRepository,
		blockRepository,
	)
	realtimeService := services.NewRealtimeService(
		validator,
		data.NotezyDB,
		blockPackRepository,
	)
	routineTagService := services.NewRoutineTagService(
		validator,
		data.NotezyDB,
		routineTagRepository,
	)
	routineTaskRecordService := services.NewRoutineTaskRecordService(
		validator,
		data.NotezyDB,
		routineTaskRecordRepository,
	)
	subShelfService := services.NewSubShelfService(
		validator,
		data.NotezyDB,
		inMemoryStorage,
		subShelfScope,
		subShelfRepository,
		rootShelfRepository,
		materialRepository,
		blockPackRepository,
	)
	blockPackService := services.NewBlockPackService(
		validator,
		data.NotezyDB,
		blockPackScope,
		subShelfRepository,
		blockPackRepository,
	)
	materialService := services.NewMaterialService(
		validator,
		data.NotezyDB,
		inMemoryStorage,
		materialScope,
		subShelfRepository,
		materialRepository,
	)
	routineService := services.NewRoutineService(
		validator,
		data.NotezyDB,
		routineScope,
		stationRepository,
		routineRepository,
		routineTagRepository,
		routineTaskRepository,
		itemRepository,
	)
	routineTaskService := services.NewRoutineTaskService(
		validator,
		data.NotezyDB,
		routineTaskScope,
		routineTaskRepository,
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
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.GET("/readyz", func(ctx *gin.Context) {
		if err := platformkafka.CheckDefaultProducer(ctx.Request.Context()); err != nil {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}

		ctx.Status(http.StatusOK)
	})

	return router
}

func Start() func() {
	shutdownObservability := observability.Initialize(context.Background())

	data.NotezyDB = data.ConnectToDatabase(config.PostgresDatabaseConfig)

	if err := userdata.Register(context.Background(), platformredis.DefaultClientManager); err != nil {
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
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(exception)
	}
	if err := platformkafka.ConnectDefaultProducer(context.Background()); err != nil {
		logs.NotezyLogger.Warn(
			context.Background(),
			"Kafka is unavailable; Core is running in degraded mode",
			attribute.String("error.message", err.Error()),
		)
	}
	outboxRelay := workers.NewOutboxRelay(
		data.NotezyDB,
		repositories.NewOutboxEventRepository(),
	)
	shutdownOutboxRelay := outboxRelay.Start(context.Background())
	yjsCommandConsumer := workers.NewYjsCommandConsumer(
		data.NotezyDB,
		services.NewYjsPersistenceService(data.NotezyDB),
		services.NewBlockService(
			validation.New(),
			data.NotezyDB,
			scopes.NewBlockScope(),
			scopes.NewBlockPackScope(),
			scopes.NewSubShelfScope(),
			repositories.NewBlockPackRepository(scopes.NewBlockPackScope()),
			repositories.NewBlockRepository(scopes.NewBlockScope()),
		),
	)
	shutdownYjsCommandConsumer := yjsCommandConsumer.Start(context.Background())

	coreTransportListener, err := net.Listen("tcp", config.CoreListenAddress())
	if err != nil {
		fmt.Println("Failed to listen for Core service transport: ", err)
		shutdownYjsCommandConsumer()
		shutdownOutboxRelay()
		platformkafka.CloseDefaultProducer()
		_ = platformredis.DefaultClientManager.DisconnectAll()
		_ = data.DisconnectToDatabase(data.NotezyDB)
		shutdownObservability()
		panic(err)
	}
	coreTransportServer := &http.Server{
		Handler: NewCoreTransportRouter(),
	}

	go func() {
		if err := coreTransportServer.Serve(coreTransportListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := coreTransportServer.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown Core service transport: ", err)
		}
		shutdownYjsCommandConsumer()
		shutdownOutboxRelay()
		platformkafka.CloseDefaultProducer()
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
		_ = data.DisconnectToDatabase(data.NotezyDB)
		shutdownObservability()
	}
}
