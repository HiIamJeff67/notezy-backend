package core

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	authcode "github.com/HiIamJeff67/notezy-backend/shared/lib/authcode"

	platformdatabase "github.com/HiIamJeff67/notezy-backend/shared/platform/database"
	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"

	coreconfig "github.com/HiIamJeff67/notezy-backend/internal/core/configs"
	userdata "github.com/HiIamJeff67/notezy-backend/internal/core/data/cache/userdata"
	data "github.com/HiIamJeff67/notezy-backend/internal/core/data/database"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/scopes"
	storage "github.com/HiIamJeff67/notezy-backend/internal/core/data/storage"
	authservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/auth"
	blockservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/blocks"
	materialservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/material"
	otherservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/other"
	realtimeservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/realtime"
	routineservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/routines"
	shelfservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/shelves"
	userservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/user"
	coretransports "github.com/HiIamJeff67/notezy-backend/internal/core/transports"
	durablejobtransport "github.com/HiIamJeff67/notezy-backend/internal/core/transports/durablejob"
	durablejobconsumers "github.com/HiIamJeff67/notezy-backend/internal/core/transports/durablejob/consumers"
	durablejobproducers "github.com/HiIamJeff67/notezy-backend/internal/core/transports/durablejob/producers"
	durablejobrouters "github.com/HiIamJeff67/notezy-backend/internal/core/transports/durablejob/routers"
	emailtransport "github.com/HiIamJeff67/notezy-backend/internal/core/transports/email"
	coremiddlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
	gatewayrouters "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/routers"
	status "github.com/HiIamJeff67/notezy-backend/internal/core/transports/status"
	yjsworkertransport "github.com/HiIamJeff67/notezy-backend/internal/core/transports/yjsworker"
	yjsworkerconsumers "github.com/HiIamJeff67/notezy-backend/internal/core/transports/yjsworker/consumers"
	yjsworkerproducers "github.com/HiIamJeff67/notezy-backend/internal/core/transports/yjsworker/producers"
	validation "github.com/HiIamJeff67/notezy-backend/internal/core/validations"
	coreworkers "github.com/HiIamJeff67/notezy-backend/internal/core/workers"
)

type Application struct {
	healthy atomic.Bool
	ready   atomic.Bool
}

func (a *Application) IsHealthy() bool {
	return a.healthy.Load()
}

func (a *Application) IsReady() bool {
	return a.ready.Load()
}

func NewCoreTransportRouter(
	config coreconfig.Config,
	kafkaProducer *platformkafka.Producer,
	userDataCacheClient *userdata.UserDataCacheClient,
) *gin.Engine {
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
	outboxEventRepository := repositories.NewOutboxEventRepository()
	inMemoryStorage := storage.NewInMemoryStorage()

	oauthService := authservices.NewOAuthService(config.OAuthGoogle.OAuthConfig())
	emailClient := emailtransport.NewClient(
		data.NotezyDB,
	)

	authService := authservices.NewAuthService(
		validator,
		data.NotezyDB,
		userRepository,
		userInfoRepository,
		userAccountRepository,
		userSettingRepository,
		rootShelfRepository,
		outboxEventRepository,
		oauthService,
		emailClient,
		userDataCacheClient,
		authcode.New(),
	)
	rootShelfService := shelfservices.NewRootShelfService(
		validator,
		data.NotezyDB,
		rootShelfScope,
		rootShelfRepository,
		usersToShelvesRepository,
		blockPackRepository,
	)
	stationService := routineservices.NewStationService(
		validator,
		data.NotezyDB,
		stationScope,
		stationRepository,
		usersToStationsRepository,
	)
	userSettingService := userservices.NewUserSettingService(
		validator,
		data.NotezyDB,
		userSettingRepository,
	)
	userInfoService := userservices.NewUserInfoService(
		validator,
		data.NotezyDB,
		userInfoRepository,
		userDataCacheClient,
	)
	userAccountService := userservices.NewUserAccountService(
		validator,
		data.NotezyDB,
		userRepository,
		userAccountRepository,
		oauthService,
	)
	userService := userservices.NewUserService(
		validator,
		data.NotezyDB,
		userRepository,
		userDataCacheClient,
	)
	blockService := blockservices.NewBlockService(
		validator,
		data.NotezyDB,
		blockScope,
		blockPackScope,
		subShelfScope,
		blockPackRepository,
		blockRepository,
	)
	realtimeService := realtimeservices.NewRealtimeService(
		validator,
		data.NotezyDB,
		blockPackRepository,
	)
	routineTagService := routineservices.NewRoutineTagService(
		validator,
		data.NotezyDB,
		routineTagRepository,
	)
	routineTaskRecordService := routineservices.NewRoutineTaskRecordService(
		validator,
		data.NotezyDB,
		routineTaskRecordRepository,
	)
	subShelfService := shelfservices.NewSubShelfService(
		validator,
		data.NotezyDB,
		inMemoryStorage,
		subShelfScope,
		subShelfRepository,
		rootShelfRepository,
		materialRepository,
		blockPackRepository,
	)
	blockPackService := blockservices.NewBlockPackService(
		validator,
		data.NotezyDB,
		blockPackScope,
		subShelfRepository,
		blockPackRepository,
	)
	materialService := materialservices.NewMaterialService(
		validator,
		data.NotezyDB,
		inMemoryStorage,
		materialScope,
		subShelfRepository,
		materialRepository,
		config.StorageKeySalt,
	)
	routineService := routineservices.NewRoutineService(
		validator,
		data.NotezyDB,
		routineScope,
		stationRepository,
		routineRepository,
		routineTagRepository,
		routineTaskRepository,
		itemRepository,
	)
	routineTaskExecutionService := routineservices.NewRoutineTaskExecutionService(
		validator,
		data.NotezyDB,
		nil,
	)
	routineTaskService := routineservices.NewRoutineTaskService(
		validator,
		data.NotezyDB,
		routineTaskScope,
		routineTaskRepository,
		routineTaskRecordRepository,
		routineTaskExecutionService,
	)
	themeService := otherservices.NewThemeService(data.NotezyDB)
	itemService := shelfservices.NewItemService(data.NotezyDB, itemScope)
	badgeService := otherservices.NewBadgeService(data.NotezyDB)
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
		userDataCacheClient,
	)
	durablejobrouters.ConfigureBlockProjectionRoutes(router, blockService)
	return router
}

func newKafkaConsumerConfig(
	connectionConfig platformkafka.ConnectionConfig,
	config coreconfig.KafkaConsumerConfig,
	clientId string,
	consumerGroup string,
) platformkafka.ConsumerConfig {
	return platformkafka.ConsumerConfig{
		ClientConfig: platformkafka.ClientConfig{
			ConnectionConfig: connectionConfig,
			ClientId:         clientId,
		},
		ConsumerGroup:       consumerGroup,
		MaximumAttempts:     config.MaximumAttempts,
		InitialRetryBackoff: config.InitialRetryBackoff,
		MaximumRetryBackoff: config.MaximumRetryBackoff,
		MaximumPollRecords:  config.MaximumPollRecords,
	}
}

func Start() func() {
	application := &Application{}
	config, err := coreconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	databaseConfig, err := platformdatabase.LoadConfig()
	if err != nil {
		panic(err)
	}
	redisConfig, err := platformredis.LoadConfig()
	if err != nil {
		panic(err)
	}
	kafkaConnectionConfig, err := platformkafka.LoadConnectionConfig()
	if err != nil {
		panic(err)
	}

	shutdownObservability := observability.Initialize(
		context.Background(),
		observability.LoadConfig("notezy-core"),
	)
	redisClientSet, err := platformredis.NewClientSet(redisConfig)
	if err != nil {
		shutdownObservability()
		panic(err)
	}

	data.NotezyDB = data.ConnectToDatabase(databaseConfig)
	if !data.MigrateEnumsToDatabase(data.NotezyDB) ||
		!data.MigrateTablesToDatabase(data.NotezyDB) ||
		!data.MigrateTriggersToDatabase(data.NotezyDB) ||
		!data.MigrateConstraintsToDatabase(data.NotezyDB) ||
		!data.SeedDefaultDataToDatabase(data.NotezyDB) {
		_ = data.DisconnectToDatabase(data.NotezyDB)
		shutdownObservability()
		panic(errors.New("failed to initialize Core database schema"))
	}

	userDataCacheStore, err := userdata.Register(context.Background(), redisClientSet)
	if err != nil {
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
		_ = redisClientSet.Close()
		shutdownObservability()
		panic(exception)
	}
	userDataCacheClient := userdata.NewUserDataCacheClient(config.UserDataCache, userDataCacheStore)
	kafkaProducer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: kafkaConnectionConfig,
		ClientId:         "notezy-core",
	})
	kafkaReady := err == nil
	if err != nil {
		logs.NotezyLogger.Warn(
			context.Background(),
			"Kafka is unavailable; Core is running in degraded mode",
			attribute.String("error.message", err.Error()),
		)
	} else if err := kafkaProducer.Ping(context.Background()); err != nil {
		kafkaReady = false
		logs.NotezyLogger.Warn(
			context.Background(),
			"Kafka is unavailable; Core is running in degraded mode",
			attribute.String("error.message", err.Error()),
		)
	}
	outboxRelay := coretransports.NewOutboxRelay(
		data.NotezyDB,
		repositories.NewOutboxEventRepository(),
		kafkaProducer,
		config.OutboxRelay,
	)
	shutdownOutboxRelay := outboxRelay.Start(context.Background())
	yjsMaintenanceReconciliationWorker := coreworkers.NewYjsMaintenanceReconciliationWorker(
		data.NotezyDB,
		repositories.NewOutboxEventRepository(),
	)
	shutdownYjsMaintenanceReconciliationWorker := yjsMaintenanceReconciliationWorker.Start(context.Background())
	routineTaskExecutionService := routineservices.NewRoutineTaskExecutionService(
		validation.New(),
		data.NotezyDB,
		nil,
	)
	routineTaskClaimConsumer := durablejobconsumers.NewDurableJobRoutineTaskClaimConsumer(
		routineservices.NewRoutineTaskService(
			validation.New(),
			data.NotezyDB,
			scopes.NewRoutineTaskScope(),
			repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope()),
			repositories.NewRoutineTaskRecordRepository(scopes.NewRoutineTaskRecordScope()),
			routineTaskExecutionService,
		),
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-durablejob-routine-task",
			durablejobtransport.RoutineTaskClaimConsumerGroup,
		),
	)
	shutdownRoutineTaskClaimConsumer := routineTaskClaimConsumer.Start(context.Background())
	routineTaskResultConsumer := durablejobconsumers.NewDurableJobRoutineTaskResultConsumer(
		routineservices.NewRoutineTaskService(
			validation.New(),
			data.NotezyDB,
			scopes.NewRoutineTaskScope(),
			repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope()),
			repositories.NewRoutineTaskRecordRepository(scopes.NewRoutineTaskRecordScope()),
			routineTaskExecutionService,
		),
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-durablejob-routine-task-result",
			durablejobtransport.RoutineTaskResultConsumerGroup,
		),
		routineTaskExecutionService,
	)
	shutdownRoutineTaskResultConsumer := routineTaskResultConsumer.Start(context.Background())
	yjsMaintenanceRequestConsumer := durablejobconsumers.NewYjsMaintenanceRequestConsumer(
		data.NotezyDB,
		yjsworkerproducers.NewYjsMaintenanceCommandProducer(kafkaProducer),
		durablejobproducers.NewYjsMaintenanceResultProducer(kafkaProducer),
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-durablejob-yjs-maintenance",
			durablejobtransport.YjsMaintenanceRequestConsumerGroup,
		),
	)
	shutdownYjsMaintenanceRequestConsumer := yjsMaintenanceRequestConsumer.Start(context.Background())
	yjsMaintenanceResultConsumer := yjsworkerconsumers.NewYjsMaintenanceResultConsumer(
		durablejobproducers.NewYjsMaintenanceResultProducer(kafkaProducer),
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-yjs-maintenance-result",
			yjsworkertransport.MaintenanceResultConsumerGroup,
		),
	)
	shutdownYjsMaintenanceResultConsumer := yjsMaintenanceResultConsumer.Start(context.Background())
	yjsCommandConsumer := yjsworkerconsumers.NewYjsCommandConsumer(
		data.NotezyDB,
		blockservices.NewYjsPersistenceService(data.NotezyDB),
		blockservices.NewBlockService(
			validation.New(),
			data.NotezyDB,
			scopes.NewBlockScope(),
			scopes.NewBlockPackScope(),
			scopes.NewSubShelfScope(),
			repositories.NewBlockPackRepository(scopes.NewBlockPackScope()),
			repositories.NewBlockRepository(scopes.NewBlockScope()),
		),
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-yjsworker",
			yjsworkertransport.CommandConsumerGroup,
		),
	)
	shutdownYjsCommandConsumer := yjsCommandConsumer.Start(context.Background())

	coreTransportListener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		fmt.Println("Failed to listen for Core service transport: ", err)
		shutdownYjsCommandConsumer()
		shutdownYjsMaintenanceResultConsumer()
		shutdownYjsMaintenanceRequestConsumer()
		shutdownYjsMaintenanceReconciliationWorker()
		shutdownRoutineTaskResultConsumer()
		shutdownRoutineTaskClaimConsumer()
		shutdownOutboxRelay()
		kafkaProducer.Close()
		_ = redisClientSet.Close()
		_ = data.DisconnectToDatabase(data.NotezyDB)
		shutdownObservability()
		panic(err)
	}
	application.healthy.Store(true)
	application.ready.Store(kafkaReady)
	coreTransportRouter := NewCoreTransportRouter(config, kafkaProducer, userDataCacheClient)
	status.ConfigureStartedRouter(coreTransportRouter, application.IsHealthy)
	status.ConfigureHealthRouter(coreTransportRouter, application.IsReady)
	coreTransportServer := &http.Server{
		Handler: coreTransportRouter,
	}

	go func() {
		if err := coreTransportServer.Serve(coreTransportListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		application.ready.Store(false)
		application.healthy.Store(false)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := coreTransportServer.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown Core service transport: ", err)
		}
		shutdownYjsCommandConsumer()
		shutdownYjsMaintenanceResultConsumer()
		shutdownYjsMaintenanceRequestConsumer()
		shutdownYjsMaintenanceReconciliationWorker()
		shutdownRoutineTaskResultConsumer()
		shutdownRoutineTaskClaimConsumer()
		shutdownOutboxRelay()
		kafkaProducer.Close()
		if err := redisClientSet.Close(); err != nil {
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
