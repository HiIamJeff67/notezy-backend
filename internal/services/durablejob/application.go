package durablejob

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	platformdatabase "github.com/HiIamJeff67/notezy-backend/internal/platform/database"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	durablejobconfig "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/config"
	routinetask "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask"
	coreconsumers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/transports/core/consumers"
	coreproducers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/transports/core/producers"
	corestrategies "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/transports/core/strategies"
)

func Start() func() {
	config, err := durablejobconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	databaseConfig, err := platformdatabase.LoadConfig()
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

	db := data.ConnectToDatabase(databaseConfig)
	kafkaProducer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: kafkaConnectionConfig,
		ClientId:         "notezy-durablejob",
	})
	if err != nil {
		_ = data.DisconnectToDatabase(db)
		shutdownObservability()
		panic(err)
	}
	if err := kafkaProducer.Ping(context.Background()); err != nil {
		kafkaProducer.Close()
		_ = data.DisconnectToDatabase(db)
		shutdownObservability()
		panic(err)
	}

	routineTaskEngine := routinetask.NewEngine(
		db,
		config,
	)
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
			ConsumerGroup:       durablejobeventscontract.DurableJobRoutineTaskConsumerGroup,
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
			ConsumerGroup:       durablejobeventscontract.DurableJobYjsMaintenanceConsumerGroup,
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
			ConsumerGroup:       durablejobeventscontract.DurableJobYjsMaintenanceConsumerGroup,
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
		_ = data.DisconnectToDatabase(db)
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
		_ = data.DisconnectToDatabase(db)
		shutdownObservability()
	}
}
