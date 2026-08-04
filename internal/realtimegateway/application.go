package realtimegateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"

	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/ratelimitrecord"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
	gatewayrouters "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/gateway/routers"
	realtimewebsocket "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/websocket"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/websocket/middlewares"
	workers "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/workers"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func Start() func() {
	shutdownObservability := observability.Initialize(context.Background())

	if err := realtimelease.Register(context.Background(), platformredis.DefaultClientManager); err != nil {
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	if err := ratelimitrecord.Register(context.Background(), platformredis.DefaultClientManager); err != nil {
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	if err := platformkafka.ConnectDefaultProducer(context.Background()); err != nil {
		logs.NotezyLogger.Warn(
			context.Background(),
			"Kafka is unavailable; RealtimeGateway is running in degraded mode",
			attribute.String("error.message", err.Error()),
		)
	}

	router := gin.Default()
	proxies := strings.Split(os.Getenv("GIN_TRUSTED_PROXIES"), ",")
	if err := router.SetTrustedProxies(proxies); err != nil {
		platformkafka.CloseDefaultProducer()
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.GET("/readyz", func(ctx *gin.Context) {
		if err := platformkafka.CheckDefaultProducer(ctx.Request.Context()); err != nil {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}

		ctx.Status(http.StatusOK)
	})
	gatewayrouters.ConfigureRoutes(
		router,
		realtimelease.NewRealtimeLeaseCacheClient(),
	)

	websocketClient := realtimewebsocket.NewWebSocketClient()
	lifecycleConsumer := workers.NewLifecycleConsumer(realtimelease.NewRealtimeLeaseCacheClient())
	shutdownLifecycleConsumer := lifecycleConsumer.Start(context.Background())
	routes := router.Group("/" + constants.RealtimeDevelopmentBaseURL)
	routes.Use(
		middlewares.DomainWhiteListMiddleware(),
		middlewares.UnauthorizedRateLimitMiddleware(),
	)
	routes.GET("", websocketClient.Handle)

	listenAddress := os.Getenv("REALTIME_GATEWAY_LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "0.0.0.0:7779"
	}

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		platformkafka.CloseDefaultProducer()
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	server := &http.Server{
		Handler: router,
	}

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		shutdownLifecycleConsumer()
		websocketClient.Shutdown()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown WebSocket server: ", err)
		}
		platformkafka.CloseDefaultProducer()
		middlewares.StopUnauthorizedRateLimiter()
		if err := platformredis.DefaultClientManager.DisconnectAll(); err != nil {
			fmt.Println("Failed to disconnect WebSocket cache servers: ", err)
		}
		shutdownObservability()
	}
}
