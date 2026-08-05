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

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	platformdatabase "github.com/HiIamJeff67/notezy-backend/internal/platform/database"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	coreconfig "github.com/HiIamJeff67/notezy-backend/internal/services/core/config"
	userdata "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/cache/userdata"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/scopes"
	storage "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/storage"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	coretransports "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports"
	durablejobconsumers "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/durablejob/consumers"
	durablejobproducers "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/durablejob/producers"
	durablejobrouters "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/durablejob/routers"
	emailtransport "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/email"
	coremiddlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
	gatewayrouters "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/routers"
	yjsworkerconsumers "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/yjsworker/consumers"
	yjsworkerproducers "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/yjsworker/producers"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
	coreworkers "github.com/HiIamJeff67/notezy-backend/internal/services/core/workers"
	authcode "github.com/HiIamJeff67/notezy-backend/shared/lib/authcode"
)

func NewCoreTransportRouter(
	config coreconfig.Config,
	kafkaProducer *platformkafka.Producer,
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
	inMemoryStorage := storage.NewInMemoryStorage()

	oauthService := services.NewOAuthService(config.OAuthGoogle.OAuthConfig())
	emailClient := emailtransport.NewClient(
		config.EmailBaseUrl,
		config.EmailClientTimeout,
	)

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
		authcode.New(),
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
		config.StorageKeySalt,
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
		if err := kafkaProducer.Ping(ctx.Request.Context()); err != nil {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}

		ctx.Status(http.StatusOK)
	})

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
	redisClientManager := platformredis.NewClientManager(redisConfig)

	data.NotezyDB = data.ConnectToDatabase(databaseConfig)

	if err := userdata.Register(context.Background(), redisClientManager); err != nil {
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
		_ = redisClientManager.DisconnectAll()
		shutdownObservability()
		panic(exception)
	}
	kafkaProducer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: kafkaConnectionConfig,
		ClientId:         "notezy-core",
	})
	if err != nil {
		logs.NotezyLogger.Warn(
			context.Background(),
			"Kafka is unavailable; Core is running in degraded mode",
			attribute.String("error.message", err.Error()),
		)
	} else if err := kafkaProducer.Ping(context.Background()); err != nil {
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
	routineTaskClaimConsumer := durablejobconsumers.NewDurableJobRoutineTaskClaimConsumer(
		services.NewRoutineTaskService(
			validation.New(),
			data.NotezyDB,
			scopes.NewRoutineTaskScope(),
			repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope()),
		),
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-durablejob-routine-task",
			durablejobeventscontract.CoreDurableJobRoutineTaskClaimConsumerGroup,
		),
	)
	shutdownRoutineTaskClaimConsumer := routineTaskClaimConsumer.Start(context.Background())
	routineTaskResultConsumer := durablejobconsumers.NewDurableJobRoutineTaskResultConsumer(
		services.NewRoutineTaskService(
			validation.New(),
			data.NotezyDB,
			scopes.NewRoutineTaskScope(),
			repositories.NewRoutineTaskRepository(scopes.NewRoutineTaskScope()),
		),
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-durablejob-routine-task-result",
			durablejobeventscontract.CoreDurableJobRoutineTaskResultConsumerGroup,
		),
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
			durablejobeventscontract.CoreDurableJobYjsMaintenanceConsumerGroup,
		),
	)
	shutdownYjsMaintenanceRequestConsumer := yjsMaintenanceRequestConsumer.Start(context.Background())
	yjsMaintenanceResultConsumer := yjsworkerconsumers.NewYjsMaintenanceResultConsumer(
		durablejobproducers.NewYjsMaintenanceResultProducer(kafkaProducer),
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-yjs-maintenance-result",
			durablejobeventscontract.CoreYjsWorkerMaintenanceResultConsumerGroup,
		),
	)
	shutdownYjsMaintenanceResultConsumer := yjsMaintenanceResultConsumer.Start(context.Background())
	yjsCommandConsumer := yjsworkerconsumers.NewYjsCommandConsumer(
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
		newKafkaConsumerConfig(
			kafkaConnectionConfig,
			config.KafkaConsumer,
			"notezy-core-yjsworker",
			"notezy-core-yjsworker-v1",
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
		_ = redisClientManager.DisconnectAll()
		_ = data.DisconnectToDatabase(data.NotezyDB)
		shutdownObservability()
		panic(err)
	}
	coreTransportServer := &http.Server{
		Handler: NewCoreTransportRouter(config, kafkaProducer),
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
		shutdownYjsMaintenanceResultConsumer()
		shutdownYjsMaintenanceRequestConsumer()
		shutdownYjsMaintenanceReconciliationWorker()
		shutdownRoutineTaskResultConsumer()
		shutdownRoutineTaskClaimConsumer()
		shutdownOutboxRelay()
		kafkaProducer.Close()
		if err := redisClientManager.DisconnectAll(); err != nil {
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
