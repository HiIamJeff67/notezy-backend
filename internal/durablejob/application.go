package durablejob

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"

	durablejobconfig "github.com/HiIamJeff67/notezy-backend/internal/durablejob/configs"
	routinetask "github.com/HiIamJeff67/notezy-backend/internal/durablejob/routinetask"
	coreconsumers "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/consumers"
	coreproducers "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/producers"
	corestrategies "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/strategies"
	realtimegatewayproducers "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/realtimegateway/producers"
	status "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/status"
)

type Application struct {
	healthy           atomic.Bool
	ready             atomic.Bool
	routineTaskEngine *routinetask.Engine
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) IsHealthy() bool {
	return a.healthy.Load()
}

func (a *Application) IsReady() bool {
	return a.ready.Load() && a.routineTaskEngine != nil && a.routineTaskEngine.IsReady()
}

func (a *Application) Start() func() {
	config, err := durablejobconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	kafkaConnectionConfig, err := platformkafka.LoadConnectionConfig()
	if err != nil {
		panic(err)
	}
	shutdownObservability := observability.Initialize(
		context.Background(),
		observability.LoadConfig("notezy-durable-job"),
	)

	// Establish the Kafka producer before constructing engines and consumers.
	kafkaProducer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: kafkaConnectionConfig,
		ClientId:         "notezy-durable-job",
	})
	if err != nil {
		shutdownObservability()
		panic(err)
	}
	if err := kafkaProducer.Ping(context.Background()); err != nil {
		kafkaProducer.Close()
		shutdownObservability()
		panic(err)
	}

	// Construct and start the durable-job workers that consume and publish tasks.
	routineTaskEngine := routinetask.NewEngine(config)
	a.routineTaskEngine = routineTaskEngine
	application := a
	routineTaskClaimProducer := coreproducers.NewRoutineTaskClaimProducer(kafkaProducer)
	routineTaskResultProducer := coreproducers.NewRoutineTaskResultProducer(kafkaProducer)
	routineTaskLifecycleProducer := realtimegatewayproducers.NewRoutineTaskLifecycleProducer(kafkaProducer)
	routineTaskEngine.SetResultPublisher(routineTaskResultProducer.Produce)
	routineTaskEngine.SetRoutineTaskRunningPublisher(
		routineTaskLifecycleProducer.ProduceRoutineTaskRunning,
	)
	routineTaskAssignmentConsumer := coreconsumers.NewRoutineTaskAssignmentConsumer(
		routineTaskEngine,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-durable-job-routine-task",
			},
			ConsumerGroup:       durablejobconfig.RoutineTaskConsumerGroup,
			MaximumAttempts:     config.KafkaConsumer.MaximumAttempts,
			InitialRetryBackoff: config.KafkaConsumer.InitialRetryBackoff,
			MaximumRetryBackoff: config.KafkaConsumer.MaximumRetryBackoff,
			MaximumPollRecords:  config.KafkaConsumer.MaximumPollRecords,
		},
	)
	shutdownRoutineTaskAssignmentConsumer := routineTaskAssignmentConsumer.Start(context.Background())
	shutdownRoutineTaskEngine := routineTaskEngine.Start(
		context.Background(),
		routineTaskClaimProducer.Produce,
	)

	yjsMaintenanceStrategy := corestrategies.NewYjsMaintenanceStrategy(config.YjsMaintenanceStrategy)
	yjsMaintenanceRequestProducer := coreproducers.NewYjsMaintenanceRequestProducer(kafkaProducer)
	yjsMaintenanceHintConsumer := coreconsumers.NewYjsMaintenanceHintConsumer(
		yjsMaintenanceRequestProducer,
		yjsMaintenanceStrategy,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-durable-job-yjs-maintenance",
			},
			ConsumerGroup:       durablejobconfig.YjsMaintenanceHintConsumerGroup,
			MaximumAttempts:     config.KafkaConsumer.MaximumAttempts,
			InitialRetryBackoff: config.KafkaConsumer.InitialRetryBackoff,
			MaximumRetryBackoff: config.KafkaConsumer.MaximumRetryBackoff,
			MaximumPollRecords:  config.KafkaConsumer.MaximumPollRecords,
		},
	)
	shutdownYjsMaintenanceHintConsumer := yjsMaintenanceHintConsumer.Start(context.Background())
	yjsMaintenanceResultConsumer := coreconsumers.NewYjsMaintenanceResultConsumer(
		yjsMaintenanceStrategy,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-durable-job-yjs-maintenance-result",
			},
			ConsumerGroup:       durablejobconfig.YjsMaintenanceResultConsumerGroup,
			MaximumAttempts:     config.KafkaConsumer.MaximumAttempts,
			InitialRetryBackoff: config.KafkaConsumer.InitialRetryBackoff,
			MaximumRetryBackoff: config.KafkaConsumer.MaximumRetryBackoff,
			MaximumPollRecords:  config.KafkaConsumer.MaximumPollRecords,
		},
	)
	shutdownYjsMaintenanceResultConsumer := yjsMaintenanceResultConsumer.Start(context.Background())

	// Expose health endpoints only after worker dependencies have started.
	mux := http.NewServeMux()
	status.ConfigureStartedRouter(mux, application.IsHealthy)
	status.ConfigureHealthRouter(mux, application.IsReady)
	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		shutdownYjsMaintenanceResultConsumer()
		shutdownYjsMaintenanceHintConsumer()
		shutdownRoutineTaskEngine()
		shutdownRoutineTaskAssignmentConsumer()
		kafkaProducer.Close()
		shutdownObservability()
		panic(err)
	}
	application.healthy.Store(true)
	application.ready.Store(application.routineTaskEngine.IsReady())
	server := &http.Server{
		Handler: mux,
	}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		// Stop HTTP traffic before stopping workers and Kafka.
		application.ready.Store(false)
		application.healthy.Store(false)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown DurableJob server: ", err)
		}
		shutdownYjsMaintenanceResultConsumer()
		shutdownYjsMaintenanceHintConsumer()
		shutdownRoutineTaskEngine()
		shutdownRoutineTaskAssignmentConsumer()
		kafkaProducer.Close()
		shutdownObservability()
	}
}

func Start() func() {
	return NewApplication().Start()
}
