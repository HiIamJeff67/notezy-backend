package durablejob

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"

	durablejobconfig "github.com/HiIamJeff67/notezy-backend/internal/durablejob/configs"
	routinetask "github.com/HiIamJeff67/notezy-backend/internal/durablejob/routinetask"
	coreconsumers "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/consumers"
	coreproducers "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/producers"
	corestrategies "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/strategies"
)

func Start() func() {
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
		observability.LoadConfig("notezy-durablejob"),
	)

	kafkaProducer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: kafkaConnectionConfig,
		ClientId:         "notezy-durablejob",
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

	routineTaskEngine := routinetask.NewEngine(config)
	routineTaskClaimProducer := coreproducers.NewRoutineTaskClaimProducer(kafkaProducer)
	routineTaskResultProducer := coreproducers.NewRoutineTaskResultProducer(kafkaProducer)
	routineTaskEngine.SetResultPublisher(routineTaskResultProducer.Produce)
	routineTaskAssignmentConsumer := coreconsumers.NewRoutineTaskAssignmentConsumer(
		routineTaskEngine,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-durablejob-routine-task",
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

	yjsMaintenanceStrategy := corestrategies.NewYjsMaintenanceStrategy()
	yjsMaintenanceRequestProducer := coreproducers.NewYjsMaintenanceRequestProducer(kafkaProducer)
	yjsMaintenanceHintConsumer := coreconsumers.NewYjsMaintenanceHintConsumer(
		yjsMaintenanceRequestProducer,
		yjsMaintenanceStrategy,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-durablejob-yjs-maintenance",
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
				ClientId:         "notezy-durablejob-yjs-maintenance-result",
			},
			ConsumerGroup:       durablejobconfig.YjsMaintenanceResultConsumerGroup,
			MaximumAttempts:     config.KafkaConsumer.MaximumAttempts,
			InitialRetryBackoff: config.KafkaConsumer.InitialRetryBackoff,
			MaximumRetryBackoff: config.KafkaConsumer.MaximumRetryBackoff,
			MaximumPollRecords:  config.KafkaConsumer.MaximumPollRecords,
		},
	)
	shutdownYjsMaintenanceResultConsumer := yjsMaintenanceResultConsumer.Start(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/readyz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
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
	server := &http.Server{
		Handler: mux,
	}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
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
