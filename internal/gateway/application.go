package gateway

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

	observability "github.com/HiIamJeff67/notezy-backend/shared/platform/observability"
	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"

	gatewayconfig "github.com/HiIamJeff67/notezy-backend/internal/gateway/configs"
	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/gateway/data/cache/ratelimitrecord"
	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/gateway/ratelimit"
	ratelimitmiddlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	developmentroutes "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/routes/developmentroutes"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	status "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/status"
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

	rateLimitRecordCacheClient := ratelimitrecord.NewRateLimitRecordCacheClient(ratelimitrecord.Config{
		ServerRange: types.Range[int, int]{
			Start: config.Redis.RateLimitRecordServerStart,
			Size:  config.Redis.RateLimitRecordServerSize,
		},
	})
	if err := ratelimitrecord.Register(context.Background(), redisClientManager, rateLimitRecordCacheClient); err != nil {
		_ = redisClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	unauthorizedRateLimitConfig := ratelimit.DefaultUnauthorizedConfig()
	unauthorizedRateLimitConfig.CacheServerRange = rateLimitRecordCacheClient.Range
	authorizedRateLimitConfig := ratelimit.DefaultAuthorizedConfig()
	authorizedRateLimitConfig.CacheServerRange = rateLimitRecordCacheClient.Range
	ratelimitmiddlewares.InitUnauthorizedRateLimiter(unauthorizedRateLimitConfig)
	ratelimitmiddlewares.InitAuthorizedRateLimiter(authorizedRateLimitConfig)
	developmentroutes.DevelopmentRouter = gin.Default()
	if err := developmentroutes.DevelopmentRouter.SetTrustedProxies(config.TrustedProxies); err != nil {
		ratelimitmiddlewares.StopUnauthorizedRateLimiter()
		ratelimitmiddlewares.StopAuthorizedRateLimiter()
		_ = redisClientManager.DisconnectAll()
		shutdownObservability()
		panic(err)
	}
	status.ConfigureStartedRouter(developmentroutes.DevelopmentRouter, application.IsHealthy)
	status.ConfigureHealthRouter(developmentroutes.DevelopmentRouter, application.IsReady)
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
		coreadapters.NewCoreAdapter(config.CoreBaseUrl, config.CoreAdapterTimeout),
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
		application.ready.Store(false)
		application.healthy.Store(false)
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
