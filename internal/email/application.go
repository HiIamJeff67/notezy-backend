package email

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"

	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"

	emailconfig "github.com/HiIamJeff67/notezy-backend/internal/email/configs"
	renderers "github.com/HiIamJeff67/notezy-backend/internal/email/renderers"
	emailsenders "github.com/HiIamJeff67/notezy-backend/internal/email/senders"
	coretransport "github.com/HiIamJeff67/notezy-backend/internal/email/transports/core"
	status "github.com/HiIamJeff67/notezy-backend/internal/email/transports/status"
)

type Application struct {
	healthy atomic.Bool
	ready   atomic.Bool
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) IsHealthy() bool {
	return a.healthy.Load()
}

func (a *Application) IsReady() bool {
	return a.ready.Load()
}

func (a *Application) Start() func() {
	application := a
	config, err := emailconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	shutdownObservability := observability.Initialize(
		context.Background(),
		observability.LoadConfig("notezy-email"),
	)

	// Initialize renderers, the bounded sender queue, and the Kafka consumer.
	deliverySender := emailsenders.NewEmailSender(config.SMTP)
	emailWorkerManager := NewEmailWorkerManager(16, deliverySender)
	welcomeRenderer, err := renderers.NewRenderer(config.Renderers.Welcome)
	if err != nil {
		emailWorkerManager.Shutdown()
		shutdownObservability()
		panic(err)
	}
	validationRenderer, err := renderers.NewRenderer(config.Renderers.Validation)
	if err != nil {
		emailWorkerManager.Shutdown()
		shutdownObservability()
		panic(err)
	}
	securityAlertRenderer, err := renderers.NewRenderer(config.Renderers.SecurityAlert)
	if err != nil {
		emailWorkerManager.Shutdown()
		shutdownObservability()
		panic(err)
	}
	sender := coretransport.NewSender(
		emailsenders.NewWelcomeEmailSender(welcomeRenderer, emailWorkerManager.Enqueue),
		emailsenders.NewValidationEmailSender(validationRenderer, emailWorkerManager.Enqueue),
		emailsenders.NewSecurityAlertEmailSender(securityAlertRenderer, emailWorkerManager.Enqueue),
	)
	validation := validator.New()
	emailRequestConsumer := coretransport.NewEmailRequestConsumer(sender, validation, config.KafkaConsumer)
	shutdownEmailRequestConsumer := emailRequestConsumer.Start(context.Background())

	// Bind the health server after the email consumer is running.
	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		shutdownEmailRequestConsumer()
		shutdownObservability()
		panic(err)
	}

	mux := http.NewServeMux()
	status.ConfigureStartedRouter(mux, application.IsHealthy)
	status.ConfigureHealthRouter(mux, application.IsReady)
	server := &http.Server{Handler: mux}
	application.healthy.Store(true)
	application.ready.Store(true)
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		// Drain email work after HTTP traffic has stopped, then release observability.
		application.ready.Store(false)
		application.healthy.Store(false)
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

func Start() func() {
	return NewApplication().Start()
}
