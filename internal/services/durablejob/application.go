package durablejob

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data"
	routinetask "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask"
	yjsmaintenance "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/yjsmaintenance"
)

func Start() func() {
	shutdownObservability := observability.Initialize(context.Background())

	db := data.ConnectToDatabase(config.PostgresDatabaseConfig)

	routineTaskEngine := routinetask.NewEngine(db)
	shutdownRoutineTaskEngine := routineTaskEngine.Start(context.Background())

	yjsMaintenanceEngine := yjsmaintenance.NewEngine(db)
	shutdownYjsMaintenanceEngine := yjsMaintenanceEngine.Start(context.Background())

	listenAddress := os.Getenv("DURABLEJOB_LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "0.0.0.0:8082"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/readyz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		shutdownYjsMaintenanceEngine()
		shutdownRoutineTaskEngine()
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
		_ = data.DisconnectToDatabase(db)
		shutdownObservability()
	}
}
