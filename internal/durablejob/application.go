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

type ApplicationInterface interface {
	loadConfig() durablejobconfig.Config
	loadKafkaConnectionConfig() platformkafka.ConnectionConfig
	initializeObservability() func()
	initializeKafka(platformkafka.ConnectionConfig, func()) *platformkafka.Producer
	initializeWorkers(durablejobconfig.Config, platformkafka.ConnectionConfig, *platformkafka.Producer) func()
	buildRouter() *http.ServeMux
	startHTTP(durablejobconfig.Config, *http.ServeMux, func(), *platformkafka.Producer, func()) func()
	Start() func()
	IsHealthy() bool
	IsReady() bool
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

func (a *Application) loadConfig() durablejobconfig.Config {
	config, err := durablejobconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	return config
}

func (a *Application) loadKafkaConnectionConfig() platformkafka.ConnectionConfig {
	kafkaConnectionConfig, err := platformkafka.LoadConnectionConfig()
	if err != nil {
		panic(err)
	}
	return kafkaConnectionConfig
}

func (a *Application) initializeObservability() func() {
	return observability.Initialize(
		context.Background(),
		observability.LoadConfig("notezy-durable-job"),
	)
}

func (a *Application) initializeKafka(
	config platformkafka.ConnectionConfig,
	shutdownObservability func(),
) *platformkafka.Producer {
	kafkaProducer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: config,
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
	return kafkaProducer
}

func (a *Application) initializeWorkers(
	config durablejobconfig.Config,
	kafkaConnection platformkafka.ConnectionConfig,
	kafkaProducer *platformkafka.Producer,
) func() {
	// Construct and start the durable-job workers that consume and publish tasks.
	routineTaskEngine := routinetask.NewEngine(config)
	a.routineTaskEngine = routineTaskEngine
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
				ConnectionConfig: kafkaConnection,
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
				ConnectionConfig: kafkaConnection,
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
				ConnectionConfig: kafkaConnection,
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

	return func() {
		shutdownYjsMaintenanceResultConsumer()
		shutdownYjsMaintenanceHintConsumer()
		shutdownRoutineTaskEngine()
		shutdownRoutineTaskAssignmentConsumer()
	}
}

func (a *Application) buildRouter() *http.ServeMux {
	mux := http.NewServeMux()
	status.ConfigureStartedRouter(mux, a.IsHealthy)
	status.ConfigureHealthRouter(mux, a.IsReady)
	return mux
}

func (a *Application) startHTTP(
	config durablejobconfig.Config,
	mux *http.ServeMux,
	shutdownWorkers func(),
	kafkaProducer *platformkafka.Producer,
	shutdownObservability func(),
) func() {
	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		shutdownWorkers()
		kafkaProducer.Close()
		shutdownObservability()
		panic(err)
	}
	a.healthy.Store(true)
	a.ready.Store(a.routineTaskEngine.IsReady())
	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		// Stop HTTP traffic before stopping workers and Kafka.
		a.ready.Store(false)
		a.healthy.Store(false)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown DurableJob server: ", err)
		}
		shutdownWorkers()
		kafkaProducer.Close()
		shutdownObservability()
	}
}

func (a *Application) Start() func() {
	shutdownObservability := a.initializeObservability()
	config := a.loadConfig()
	kafkaConnection := a.loadKafkaConnectionConfig()
	kafkaProducer := a.initializeKafka(kafkaConnection, shutdownObservability)
	shutdownWorkers := a.initializeWorkers(config, kafkaConnection, kafkaProducer)
	router := a.buildRouter()
	return a.startHTTP(config, router, shutdownWorkers, kafkaProducer, shutdownObservability)
}

// make sure Application struct followed the ApplicationInterface implementations
var _ ApplicationInterface = (*Application)(nil)
