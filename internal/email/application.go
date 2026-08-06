package email

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"

	emailconfig "github.com/HiIamJeff67/notezy-backend/internal/email/configs"
	"github.com/HiIamJeff67/notezy-backend/internal/email/renderers"
	emailsenders "github.com/HiIamJeff67/notezy-backend/internal/email/senders"
	coretransport "github.com/HiIamJeff67/notezy-backend/internal/email/transports/core"
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
	deliverySender := emailsenders.NewEmailSender(config.SMTP)
	emailWorkerManager := NewEmailWorkerManager(16, deliverySender)
	welcomeRenderer, exception := renderers.NewRenderer(config.Renderers.Welcome)
	if exception != nil {
		emailWorkerManager.Shutdown()
		shutdownObservability()
		panic(exception)
	}
	validationRenderer, exception := renderers.NewRenderer(config.Renderers.Validation)
	if exception != nil {
		emailWorkerManager.Shutdown()
		shutdownObservability()
		panic(exception)
	}
	securityAlertRenderer, exception := renderers.NewRenderer(config.Renderers.SecurityAlert)
	if exception != nil {
		emailWorkerManager.Shutdown()
		shutdownObservability()
		panic(exception)
	}
	sender := coretransport.NewSender(
		emailsenders.NewWelcomeEmailSender(welcomeRenderer, emailWorkerManager.Enqueue),
		emailsenders.NewValidationEmailSender(validationRenderer, emailWorkerManager.Enqueue),
		emailsenders.NewSecurityAlertEmailSender(securityAlertRenderer, emailWorkerManager.Enqueue),
	)
	validation := validator.New()
	emailRequestConsumer := coretransport.NewEmailRequestConsumer(sender, validation, config.KafkaConsumer)
	shutdownEmailRequestConsumer := emailRequestConsumer.Start(context.Background())

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		shutdownEmailRequestConsumer()
		shutdownObservability()
		panic(err)
	}

	server := &http.Server{
		Handler: coretransport.NewRouter(),
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
		shutdownEmailRequestConsumer()
		emailWorkerManager.Shutdown()
		shutdownObservability()
	}
}
