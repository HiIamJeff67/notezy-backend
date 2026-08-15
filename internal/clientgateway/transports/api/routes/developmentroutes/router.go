package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/ratelimit"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
	notificationadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/notification/adapters"
)

var (
	DevelopmentRouter         *gin.Engine
	DevelopmentAPIRouterGroup *gin.RouterGroup
)

type RateLimiters struct {
	Unauthorized *ratelimit.HybridRateLimiter
	Authorized   *ratelimit.HybridRateLimiter
}

type APIRouteDependencies struct {
	CoreClient                *coreadapters.CoreAdapter
	NotificationClient        *notificationadapters.NotificationAdapter
	AllowedDomains            []string
	AccessTokenCookieHandler  *cookies.CookieHandler
	RefreshTokenCookieHandler *cookies.CookieHandler
	RateLimiters              RateLimiters
}

func NewRouter(deps APIRouteDependencies) *gin.Engine {
	DevelopmentRouter = gin.Default()
	coreClient, notificationClient := deps.CoreClient, deps.NotificationClient
	allowedDomains, accessTokenCookieHandler := deps.AllowedDomains, deps.AccessTokenCookieHandler
	refreshTokenCookieHandler, rateLimiters := deps.RefreshTokenCookieHandler, deps.RateLimiters
	DevelopmentAPIRouterGroup = DevelopmentRouter.Group("/" + gatewaycontract.APIDevelopmentBaseURL) // use in development mode
	DevelopmentAPIRouterGroup.Use(
		middlewares.SanitizeXForwardedForMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.DomainWhiteListMiddleware(allowedDomains),
	)
	DevelopmentAPIRouterGroup.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(200) })
	fmt.Println("API router group path:", DevelopmentAPIRouterGroup.BasePath())

	configureDevelopmentAuthRoutes(DevelopmentAPIRouterGroup, AuthRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentUserRoutes(DevelopmentAPIRouterGroup, UserRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentUserInfoRoutes(DevelopmentAPIRouterGroup, UserInfoRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureUserSettingRoutes(DevelopmentAPIRouterGroup, UserSettingRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentUserAccountRoutes(DevelopmentAPIRouterGroup, UserAccountRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentAPIKeyRoutes(DevelopmentAPIRouterGroup, APIKeyRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})

	configureDevelopmentStationRoutes(DevelopmentAPIRouterGroup, StationRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentRoutineRoutes(DevelopmentAPIRouterGroup, RoutineRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentRoutineTagRoutes(DevelopmentAPIRouterGroup, RoutineTagRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentRoutineTaskRoutes(DevelopmentAPIRouterGroup, RoutineTaskRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentRootShelfRoutes(DevelopmentAPIRouterGroup, RootShelfRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentSubShelfRoutes(DevelopmentAPIRouterGroup, SubShelfRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentMaterialRoutes(DevelopmentAPIRouterGroup, MaterialRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentBlockPackRoutes(DevelopmentAPIRouterGroup, BlockPackRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentBlockRoutes(DevelopmentAPIRouterGroup, BlockRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})

	configureDevelopmentRoutineTaskRecordRoutes(DevelopmentAPIRouterGroup, RoutineTaskRecordRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentRealtimeRoutes(DevelopmentAPIRouterGroup, RealtimeRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentGraphQLRoutes(DevelopmentAPIRouterGroup, GraphQLRouteDependencies{CoreClient: coreClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})
	configureDevelopmentNotificationRoutes(DevelopmentAPIRouterGroup, NotificationRouteDependencies{NotificationClient: notificationClient, AccessTokenCookieHandler: accessTokenCookieHandler, RefreshTokenCookieHandler: refreshTokenCookieHandler, RateLimiters: rateLimiters})

	configureStaticRoutes(DevelopmentAPIRouterGroup, StaticRouteDependencies{RateLimiters: rateLimiters})

	return DevelopmentRouter
}
