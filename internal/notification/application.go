package notification

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"
	sharedvalidations "github.com/HiIamJeff67/notezy-backend/shared/validations"

	configs "github.com/HiIamJeff67/notezy-backend/internal/notification/configs"
	database "github.com/HiIamJeff67/notezy-backend/internal/notification/data/database"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/notification/data/database/repositories"
	services "github.com/HiIamJeff67/notezy-backend/internal/notification/services"
	consumers "github.com/HiIamJeff67/notezy-backend/internal/notification/transports/core/consumers"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/notification/transports/gateway/endpoints"
	routers "github.com/HiIamJeff67/notezy-backend/internal/notification/transports/gateway/routers"
	validations "github.com/HiIamJeff67/notezy-backend/internal/notification/validations"
	workers "github.com/HiIamJeff67/notezy-backend/internal/notification/workers"
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

func Start() func() {
	application := &Application{}
	config, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}
	shutdownObservability := observability.Initialize(
		context.Background(),
		observability.LoadConfig("notezy-notification"),
	)
	db, err := database.Connect(config.Database)
	if err != nil {
		shutdownObservability()
		panic(err)
	}
	producer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: config.Kafka.Connection,
		ClientId:         "notezy-notification-producer",
	})
	if err != nil {
		_ = database.Disconnect(db)
		shutdownObservability()
		panic(err)
	}
	repository := repositories.NewNotificationRepository(db)
	notificationValidator := validator.New()
	sharedvalidations.RegisterStringsValidation(notificationValidator)
	sharedvalidations.RegisterTimesValidation(notificationValidator)
	validations.RegisterNotificationValidation(notificationValidator)
	validations.RegisterNewsValidation(notificationValidator)
	validations.RegisterWarningValidation(notificationValidator)
	validations.RegisterImportantValidation(notificationValidator)
	service := services.NewNotificationService(repository, notificationValidator)
	consumer := consumers.NewNotificationRequestConsumer(
		service,
		config.Kafka.ConsumerConfig(),
	)
	relay := workers.NewOutboxRelay(
		repository,
		producer,
		config.OutboxPollInterval,
		config.OutboxClaimTimeout,
		config.OutboxInitialBackoff,
		config.OutboxMaximumBackoff,
		config.OutboxBatchSize,
		config.OutboxCleanupInterval,
		config.OutboxRetention,
	)
	cleanup := workers.NewCleanupWorker(
		service,
		config.OutboxCleanupInterval,
		config.NotificationRetention,
	)
	shutdownConsumer := consumer.Start(context.Background())
	shutdownRelay := relay.Start(context.Background())
	shutdownCleanup := cleanup.Start(context.Background())

	router := gin.New()
	router.GET("/healthz", func(ctx *gin.Context) {
		if !application.IsReady() {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}
		ctx.Status(http.StatusNoContent)
	})
	router.GET("/startedz", func(ctx *gin.Context) {
		if !application.IsHealthy() {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}
		ctx.Status(http.StatusNoContent)
	})
	endpoint := endpoints.NewNotificationEndpoint(service)
	routers.ConfigureNotificationRoutes(router.Group("/internal/v1"), endpoint)

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		shutdownCleanup()
		shutdownRelay()
		shutdownConsumer()
		producer.Close()
		_ = database.Disconnect(db)
		shutdownObservability()
		panic(err)
	}
	application.healthy.Store(true)
	application.ready.Store(true)
	server := &http.Server{Handler: router}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		application.ready.Store(false)
		application.healthy.Store(false)
		shutdownCleanup()
		shutdownRelay()
		shutdownConsumer()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			fmt.Printf("Failed to shutdown Notification server: %v\n", err)
		}
		producer.Close()
		if err := database.Disconnect(db); err != nil {
			fmt.Printf("Failed to disconnect Notification database: %v\n", err)
		}
		shutdownObservability()
	}
}
