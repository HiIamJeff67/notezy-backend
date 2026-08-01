package email

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	coretransport "github.com/HiIamJeff67/notezy-backend/internal/services/email/transports/core"
)

func Start() {
	shutdownObservability := observability.Initialize(context.Background())
	defer shutdownObservability()
	defer NotezyEmailWorkerManager.Shutdown()

	listenAddress := os.Getenv("EMAIL_LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "127.0.0.1:8081"
	}
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler: coretransport.NewRouter(coretransport.Sender{
			SendWelcomeEmail:       AsyncSendWelcomeEmail,
			SendValidationEmail:    AsyncSendValidationEmail,
			SendSecurityAlertEmail: AsyncSendSecurityAlertEmail,
		}),
	}
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			panic(err)
		}
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}
}
