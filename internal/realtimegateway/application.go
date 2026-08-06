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

	types "github.com/HiIamJeff67/notezy-backend/shared/types"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtime-gateway/v1"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"
	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"

	realtimeconfig "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/configs"
	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/ratelimitrecord"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/ratelimit"
	gatewayrouters "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/gateway/routers"
	status "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/status"
	yjsworker "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/yjsworker"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/yjsworker/middlewares"
	workers "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/workers"
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
	redisClientManager := platformredis.NewClientManager(redisConfig)
	realtimeLeaseCacheClient := realtimelease.NewRealtimeLeaseCacheClient(realtimelease.Config{
		ServerRange: types.Range[int, int]{
			Start: config.Redis.RealtimeLeaseServerStart,
			Size:  config.Redis.RealtimeLeaseServerSize,
		},
	})
	rateLimitRecordCacheClient := ratelimitrecord.NewRateLimitRecordCacheClient(ratelimitrecord.Config{
		ServerRange: types.Range[int, int]{
			Start: config.Redis.RateLimitRecordServerStart,
			Size:  config.Redis.RateLimitRecordServerSize,
		},
	})

	if err := realtimelease.Register(context.Background(), redisClientManager, realtimeLeaseCacheClient); err != nil {
		_ = redisClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	if err := ratelimitrecord.Register(context.Background(), redisClientManager, rateLimitRecordCacheClient); err != nil {
		_ = redisClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	upgradeRateLimitConfig := ratelimit.DefaultUpgradeConfig()
	upgradeRateLimitConfig.CacheServerRange = rateLimitRecordCacheClient.Range
	middlewares.InitUnauthorizedRateLimiter(upgradeRateLimitConfig)

	router := gin.Default()
	if err := router.SetTrustedProxies(config.TrustedProxies); err != nil {
		middlewares.StopUnauthorizedRateLimiter()
		_ = redisClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	status.ConfigureStartedRouter(router, application.IsHealthy)
	status.ConfigureHealthRouter(router, application.IsReady)
	gatewayrouters.ConfigureRoutes(
		router,
		realtimeLeaseCacheClient,
	)

	websocketClient := yjsworker.NewWebSocketClient(config, realtimeLeaseCacheClient)
	lifecycleConsumer := workers.NewLifecycleConsumer(
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
	routes := router.Group("/" + realtimegatewaycontract.RealtimeDevelopmentBaseURL)
	routes.Use(
		middlewares.DomainWhiteListMiddleware(config.AllowedDomains),
		middlewares.UnauthorizedRateLimitMiddleware(),
	)
	routes.GET("", websocketClient.Handle)

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		_ = redisClientManager.DisconnectAll()
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
		websocketClient.Shutdown()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown WebSocket server: ", err)
		}
		middlewares.StopUnauthorizedRateLimiter()
		if err := redisClientManager.DisconnectAll(); err != nil {
			fmt.Println("Failed to disconnect WebSocket cache servers: ", err)
		}
		shutdownObservability()
	}
}
