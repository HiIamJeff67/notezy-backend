package apigateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"

	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"
	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"

	gatewayconfig "github.com/HiIamJeff67/notezy-backend/internal/apigateway/configs"
	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/apigateway/data/cache/ratelimitrecord"
	ratelimitmiddlewares "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/middlewares"
	developmentroutes "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/routes/developmentroutes"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/core/adapters"
	status "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/status"
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
	config, err := gatewayconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	redisConfig, err := platformredis.LoadConfig()
	if err != nil {
		panic(err)
	}
	shutdownObservability := observability.Initialize(
		context.Background(),
		observability.LoadConfig("notezy-api-gateway"),
	)
	redisClientSet, err := platformredis.NewClientSet(redisConfig)
	if err != nil {
		shutdownObservability()
		panic(err)
	}
	// Initialize APIGateway's Redis-backed rate-limit record cache.
	rateLimitRecordCacheStore := ratelimitrecord.Register(context.Background(), redisClientSet)
	if err := rateLimitRecordCacheStore.Initialize(context.Background()); err != nil {
		_ = redisClientSet.Close()
		shutdownObservability()
		panic(err)
	}
	rateLimitRecordCacheClient := ratelimitrecord.NewRateLimitRecordCacheClient(rateLimitRecordCacheStore)

	// Create the application-owned rate limiters. Route registration receives
	// these instances explicitly so middleware does not keep global state.
	unauthorizedRateLimitConfig := gatewayconfig.DefaultUnauthorizedRateLimitConfig()
	unauthorizedRateLimitConfig.CacheClient = rateLimitRecordCacheClient
	unauthorizedRateLimiter := ratelimitmiddlewares.InitUnauthorizedRateLimiter(unauthorizedRateLimitConfig)

	// Build the HTTP router and apply process-wide proxy and health settings.
	developmentroutes.DevelopmentRouter = gin.Default()
	if err := developmentroutes.DevelopmentRouter.SetTrustedProxies(config.TrustedProxies); err != nil {
		unauthorizedRateLimiter.Stop()
		_ = redisClientSet.Close()
		shutdownObservability()
		panic(err)
	}
	status.ConfigureStartedRouter(developmentroutes.DevelopmentRouter, application.IsHealthy)
	status.ConfigureHealthRouter(developmentroutes.DevelopmentRouter, application.IsReady)
	developmentroutes.ConfigureAPIRoutes(
		coreadapters.NewCoreAdapter(config.CoreBaseUrl, config.CoreAdapterTimeout),
		config.AllowedDomains,
		developmentroutes.RateLimiters{
			Unauthorized: unauthorizedRateLimiter,
		},
	)

	// Bind the listener only after all dependencies and routes are ready.
	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		unauthorizedRateLimiter.Stop()
		_ = redisClientSet.Close()
		shutdownObservability()
		panic(err)
	}
	application.healthy.Store(true)
	application.ready.Store(true)
	server := &http.Server{
		Handler: developmentroutes.DevelopmentRouter,
	}

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		// Shut down request handling before releasing its shared dependencies.
		application.ready.Store(false)
		application.healthy.Store(false)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown Gateway server: ", err)
		}
		unauthorizedRateLimiter.Stop()
		if err := redisClientSet.Close(); err != nil {
			fmt.Println("Failed to disconnect Gateway cache servers: ", err)
		}
		shutdownObservability()
	}
}

func Start() func() {
	return NewApplication().Start()
}
