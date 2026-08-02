package gateway

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

	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/gateway/data/cache/ratelimitrecord"
	developmentroutes "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/routes/developmentroutes"
	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
)

func Start() func() {
	shutdownObservability := observability.Initialize(context.Background())

	if err := ratelimitrecord.Register(context.Background(), platformredis.DefaultClientManager); err != nil {
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	if err := realtimelease.Register(context.Background(), platformredis.DefaultClientManager); err != nil {
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}

	developmentroutes.DevelopmentRouter = gin.Default()
	proxies := strings.Split(os.Getenv("GIN_TRUSTED_PROXIES"), ",")
	if err := developmentroutes.DevelopmentRouter.SetTrustedProxies(proxies); err != nil {
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	developmentroutes.DevelopmentRouter.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	developmentroutes.DevelopmentRouter.GET("/readyz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	developmentroutes.ConfigureAPIRoutes()

	listener, err := net.Listen("tcp", config.GatewayListenAddress())
	if err != nil {
		_ = platformredis.DefaultClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	server := &http.Server{
		Handler: developmentroutes.DevelopmentRouter,
	}

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown Gateway server: ", err)
		}
		if err := platformredis.DefaultClientManager.DisconnectAll(); err != nil {
			fmt.Println("Failed to disconnect Gateway cache servers: ", err)
		}
		shutdownObservability()
	}
}
