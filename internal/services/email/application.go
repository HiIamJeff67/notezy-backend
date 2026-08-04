package email

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"

	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	coretransport "github.com/HiIamJeff67/notezy-backend/internal/services/email/transports/core"
)

func Start() func() {
	shutdownObservability := observability.Initialize(context.Background())

	listenAddress := os.Getenv("EMAIL_LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "127.0.0.1:8081"
	}
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		shutdownObservability()
		panic(err)
	}

	server := &http.Server{
		Handler: coretransport.NewRouter(
			coretransport.Sender{
				SendWelcomeEmail:       AsyncSendWelcomeEmail,
				SendValidationEmail:    AsyncSendValidationEmail,
				SendSecurityAlertEmail: AsyncSendSecurityAlertEmail,
			},
			validator.New(),
		),
	}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			panic(err)
		}
		NotezyEmailWorkerManager.Shutdown()
		shutdownObservability()
	}
}
