package realtimegateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
	platform "github.com/HiIamJeff67/notezy-backend/shared/platform"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtime-gateway/v1"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"
	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"

	realtimeconfig "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/configs"
	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/ratelimitrecord"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/ratelimit"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/api/middlewares"
	apiRouters "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/api/routers"
	coreconsumers "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/core/consumers"
	durablejobconsumers "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/durablejob/consumers"
	notificationconsumers "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/notification/consumers"
	status "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/status"
	websockettransport "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/websocket"
)

type Application struct {
	healthy atomic.Bool
	ready   atomic.Bool
}

func (a *Application) IsHealthy() bool {
	return a.healthy.Load()
}

func (a *Application) IsReady() bool {
	return a.ready.Load()
}

func Start() func() {
	application := &Application{}
	config, err := realtimeconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	redisConfig, err := platformredis.LoadConfig()
	if err != nil {
		panic(err)
	}
	kafkaConnectionConfig, err := platformkafka.LoadConnectionConfig()
	if err != nil {
		panic(err)
	}
	shutdownObservability := observability.Initialize(
		context.Background(),
		observability.LoadConfig("notezy-realtime-gateway"),
	)
	redisClientSet, err := platformredis.NewClientSet(redisConfig)
	if err != nil {
		shutdownObservability()
		panic(err)
	}
	realtimeLeaseCacheStore := realtimelease.Register(context.Background(), redisClientSet)
	if err := realtimeLeaseCacheStore.Initialize(context.Background()); err != nil {
		_ = redisClientSet.Close()
		shutdownObservability()
		panic(err)
	}
	realtimeLeaseCacheClient := realtimelease.NewRealtimeLeaseCacheClient(realtimeLeaseCacheStore)
	rateLimitRecordCacheStore := ratelimitrecord.Register(context.Background(), redisClientSet)
	if err := rateLimitRecordCacheStore.Initialize(context.Background()); err != nil {
		_ = redisClientSet.Close()
		shutdownObservability()
		panic(err)
	}
	rateLimitRecordCacheClient := ratelimitrecord.NewRateLimitRecordCacheClient(rateLimitRecordCacheStore)
	upgradeRateLimitConfig := realtimeconfig.DefaultUpgradeRateLimitConfig()
	upgradeRateLimitConfig.CacheClient = rateLimitRecordCacheClient
	unauthorizedRateLimiter := ratelimit.NewHybridRateLimiter(upgradeRateLimitConfig, false)
	authorizedRateLimiter := ratelimit.NewHybridRateLimiter(upgradeRateLimitConfig, true)

	router := gin.Default()
	if err := router.SetTrustedProxies(config.TrustedProxies); err != nil {
		unauthorizedRateLimiter.Stop()
		authorizedRateLimiter.Stop()
		_ = redisClientSet.Close()
		shutdownObservability()
		panic(err)
	}
	status.ConfigureStartedRouter(router, application.IsHealthy)
	status.ConfigureHealthRouter(router, application.IsReady)
	routes := router.Group("/" + realtimegatewaycontract.RealtimeDevelopmentBaseURL)
	routes.Use(
		middlewares.SanitizeXForwardedForMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.DomainWhiteListMiddleware(config.AllowedDomains),
		middlewares.UnauthorizedRateLimitMiddleware(unauthorizedRateLimiter),
	)
	routes.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })
	accessTokenCookieHandler := cookies.New(cookies.Config{
		Name:     cookies.ValidCookieName_AccessToken,
		Path:     "/",
		Duration: 30 * time.Minute,
		Secure:   platform.CurrentEnvironment == types.Environment_Production,
		HTTPOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	refreshTokenCookieHandler := cookies.New(cookies.Config{
		Name:     cookies.ValidCookieName_RefreshToken,
		Path:     "/",
		Duration: 14 * 24 * time.Hour,
		Secure:   platform.CurrentEnvironment == types.Environment_Production,
		HTTPOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	apiRouters.ConfigureRoutes(routes, realtimeLeaseCacheClient, accessTokenCookieHandler, refreshTokenCookieHandler, authorizedRateLimiter)

	websocketAdapter := websockettransport.NewWebSocketAdapter(config, realtimeLeaseCacheClient)
	lifecycleConsumer := coreconsumers.NewLifecycleConsumer(
		realtimeLeaseCacheClient,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-realtime-gateway-lifecycle",
			},
			ConsumerGroup:       "notezy-realtime-gateway-lifecycle-v1",
			MaximumAttempts:     config.KafkaConsumer.MaximumAttempts,
			InitialRetryBackoff: config.KafkaConsumer.InitialRetryBackoff,
			MaximumRetryBackoff: config.KafkaConsumer.MaximumRetryBackoff,
			MaximumPollRecords:  config.KafkaConsumer.MaximumPollRecords,
		},
	)
	shutdownLifecycleConsumer := lifecycleConsumer.Start(context.Background())
	routineTaskLifecycleConsumer := durablejobconsumers.NewRoutineTaskLifecycleConsumer(
		realtimeLeaseCacheClient,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-realtime-gateway-durable-job-routine-task-lifecycle",
			},
			ConsumerGroup:       "notezy-realtime-gateway-durable-job-routine-task-lifecycle-v1",
			MaximumAttempts:     config.KafkaConsumer.MaximumAttempts,
			InitialRetryBackoff: config.KafkaConsumer.InitialRetryBackoff,
			MaximumRetryBackoff: config.KafkaConsumer.MaximumRetryBackoff,
			MaximumPollRecords:  config.KafkaConsumer.MaximumPollRecords,
		},
	)
	shutdownRoutineTaskLifecycleConsumer := routineTaskLifecycleConsumer.Start(context.Background())
	notificationConsumer := notificationconsumers.NewNotificationConsumer(
		realtimeLeaseCacheClient,
		platformkafka.ConsumerConfig{
			ClientConfig: platformkafka.ClientConfig{
				ConnectionConfig: kafkaConnectionConfig,
				ClientId:         "notezy-realtime-gateway-notification",
			},
			ConsumerGroup:       "notezy-realtime-gateway-notification-v1",
			MaximumAttempts:     config.KafkaConsumer.MaximumAttempts,
			InitialRetryBackoff: config.KafkaConsumer.InitialRetryBackoff,
			MaximumRetryBackoff: config.KafkaConsumer.MaximumRetryBackoff,
			MaximumPollRecords:  config.KafkaConsumer.MaximumPollRecords,
		},
	)
	shutdownNotificationConsumer := notificationConsumer.Start(context.Background())
	routes.GET("", websocketAdapter.Handle)

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		unauthorizedRateLimiter.Stop()
		authorizedRateLimiter.Stop()
		_ = redisClientSet.Close()
		shutdownObservability()
		panic(err)
	}
	application.healthy.Store(true)
	application.ready.Store(true)
	server := &http.Server{
		Handler: router,
	}

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		application.ready.Store(false)
		application.healthy.Store(false)
		shutdownLifecycleConsumer()
		shutdownRoutineTaskLifecycleConsumer()
		shutdownNotificationConsumer()
		websocketAdapter.Shutdown()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown WebSocket server: ", err)
		}
		unauthorizedRateLimiter.Stop()
		authorizedRateLimiter.Stop()
		if err := redisClientSet.Close(); err != nil {
			fmt.Println("Failed to disconnect WebSocket cache servers: ", err)
		}
		shutdownObservability()
	}
}
