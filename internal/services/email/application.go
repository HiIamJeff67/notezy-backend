package email

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	emailconfig "github.com/HiIamJeff67/notezy-backend/internal/services/email/config"
	coretransport "github.com/HiIamJeff67/notezy-backend/internal/services/email/transports/core"
)

func Start() func() {
	config, err := emailconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	shutdownObservability := observability.Initialize(
		context.Background(),
		observability.LoadConfig("notezy-email"),
	)
	Initialize(config.SMTP)

	listener, err := net.Listen("tcp", config.ListenAddress)
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
