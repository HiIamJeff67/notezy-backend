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

func ConfigureAPIRoutes(
	coreClient *coreadapters.CoreAdapter,
	notificationClient *notificationadapters.NotificationAdapter,
	allowedDomains []string,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
	rateLimiters RateLimiters,
) {
	DevelopmentAPIRouterGroup = DevelopmentRouter.Group("/" + gatewaycontract.APIDevelopmentBaseURL) // use in development mode
	DevelopmentAPIRouterGroup.Use(
		middlewares.SanitizeXForwardedForMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.DomainWhiteListMiddleware(allowedDomains),
	)
	DevelopmentAPIRouterGroup.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(200) })
	fmt.Println("API router group path:", DevelopmentAPIRouterGroup.BasePath())

	configureDevelopmentAuthRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentUserRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentUserInfoRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureUserSettingRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentUserAccountRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)

	configureDevelopmentStationRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentRoutineRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentRoutineTagRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentRoutineTaskRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentRootShelfRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentSubShelfRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentMaterialRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentBlockPackRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentBlockRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)

	configureDevelopmentRoutineTaskRecordRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentRealtimeRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentGraphQLRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)
	configureDevelopmentNotificationRoutes(DevelopmentAPIRouterGroup, notificationClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters)

	configureStaticRoutes(DevelopmentAPIRouterGroup, rateLimiters)
}
