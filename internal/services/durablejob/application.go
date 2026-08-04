package durablejob

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	platformdatabase "github.com/HiIamJeff67/notezy-backend/internal/platform/database"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database"
	durablejobconfig "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/config"
	routinetask "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask"
	yjsmaintenance "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/yjsmaintenance"
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
	if err := platformkafka.ConnectDefaultProducer(
		context.Background(),
		platformkafka.ClientConfig{
			ConnectionConfig: kafkaConnectionConfig,
			ClientId:         "notezy-durablejob",
		},
	); err != nil {
		_ = data.DisconnectToDatabase(db)
		shutdownObservability()
		panic(err)
	}

	routineTaskEngine := routinetask.NewEngine(
		db,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-durablejob-routine-task",
			},
			ConsumerGroup:       "notezy-durablejob-routine-task-v1",
			MaximumAttempts:     config.KafkaConsumer.MaximumAttempts,
			InitialRetryBackoff: config.KafkaConsumer.InitialRetryBackoff,
			MaximumRetryBackoff: config.KafkaConsumer.MaximumRetryBackoff,
			MaximumPollRecords:  config.KafkaConsumer.MaximumPollRecords,
		},
		config,
	)
	shutdownRoutineTaskEngine := routineTaskEngine.Start(context.Background())

	yjsMaintenanceEngine := yjsmaintenance.NewEngine(db, config)
	shutdownYjsMaintenanceEngine := yjsMaintenanceEngine.Start(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/readyz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		shutdownYjsMaintenanceEngine()
		shutdownRoutineTaskEngine()
		platformkafka.CloseDefaultProducer()
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
		shutdownYjsMaintenanceEngine()
		shutdownRoutineTaskEngine()
		platformkafka.CloseDefaultProducer()
		_ = data.DisconnectToDatabase(db)
		shutdownObservability()
	}
}
