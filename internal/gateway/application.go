package gateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gatewayconfig "github.com/HiIamJeff67/notezy-backend/internal/gateway/config"
	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/gateway/data/cache/ratelimitrecord"
	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/gateway/ratelimit"
	ratelimitmiddlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	developmentroutes "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/routes/developmentroutes"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	realtimegatewayadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/realtimegateway/adapters"
	platform "github.com/HiIamJeff67/notezy-backend/internal/platform"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

func Start() func() {
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
		observability.LoadConfig("notezy-gateway"),
	)
	redisClientManager := platformredis.NewClientManager(redisConfig)

	if err := ratelimitrecord.Register(context.Background(), redisClientManager); err != nil {
		_ = redisClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	ratelimitmiddlewares.InitUnauthorizedRateLimiter(ratelimit.DefaultUnauthorizedConfig())
	ratelimitmiddlewares.InitAuthorizedRateLimiter(ratelimit.DefaultAuthorizedConfig())
	developmentroutes.DevelopmentRouter = gin.Default()
	if err := developmentroutes.DevelopmentRouter.SetTrustedProxies(config.TrustedProxies); err != nil {
		ratelimitmiddlewares.StopUnauthorizedRateLimiter()
		ratelimitmiddlewares.StopAuthorizedRateLimiter()
		_ = redisClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	developmentroutes.DevelopmentRouter.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	developmentroutes.DevelopmentRouter.GET("/readyz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	accessTokenCookieHandler := cookies.New(cookies.Config{
		Name:     cookies.ValidCookieName_AccessToken,
		Path:     "/",
		Duration: 30 * time.Minute, // 30 minutes
		Secure:   platform.CurrentEnvironment == types.Environment_Production,
		HTTPOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	refreshTokenCookieHandler := cookies.New(cookies.Config{
		Name:     cookies.ValidCookieName_RefreshToken,
		Path:     "/",
		Duration: 14 * 24 * time.Hour, // 14 days
		Secure:   platform.CurrentEnvironment == types.Environment_Production,
		HTTPOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	developmentroutes.ConfigureAPIRoutes(
		coreadapters.NewCoreClient(config.CoreBaseUrl, config.CoreClientTimeout),
		realtimegatewayadapters.NewRealtimeGatewayClient(
			config.RealtimeGatewayBaseUrl,
			config.RealtimeGatewayClientTimeout,
		),
		config.AllowedDomains,
		accessTokenCookieHandler,
		refreshTokenCookieHandler,
	)

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		ratelimitmiddlewares.StopUnauthorizedRateLimiter()
		ratelimitmiddlewares.StopAuthorizedRateLimiter()
		_ = redisClientManager.DisconnectAll()
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
		if err := redisClientManager.DisconnectAll(); err != nil {
			fmt.Println("Failed to disconnect Gateway cache servers: ", err)
		}
		ratelimitmiddlewares.StopUnauthorizedRateLimiter()
		ratelimitmiddlewares.StopAuthorizedRateLimiter()
		shutdownObservability()
	}
}
