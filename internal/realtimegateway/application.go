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

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/ratelimitrecord"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
	realtimecore "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/core"
	realtimewebsocket "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/websocket"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/websocket/middlewares"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
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

	router := gin.Default()
	proxies := strings.Split(os.Getenv("GIN_TRUSTED_PROXIES"), ",")
	if err := router.SetTrustedProxies(proxies); err != nil {
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.GET("/readyz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	websocketClient := realtimewebsocket.NewWebSocketClient(realtimecore.NewCoreClient(coreadapters.NewConfiguredCoreClient()))
	routes := router.Group("/" + constants.RealtimeDevelopmentBaseURL)
	routes.Use(
		middlewares.DomainWhiteListMiddleware(),
		middlewares.UnauthorizedRateLimitMiddleware(),
	)
	routes.GET("", websocketClient.Handle)

	listener, err := net.Listen("tcp", config.RealtimeGatewayListenAddress())
	if err != nil {
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
		websocketClient.Shutdown()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown WebSocket server: ", err)
		}
		middlewares.StopUnauthorizedRateLimiter()
		if err := platformredis.DefaultClientManager.DisconnectAll(); err != nil {
			fmt.Println("Failed to disconnect WebSocket cache servers: ", err)
		}
		shutdownObservability()
	}
}
