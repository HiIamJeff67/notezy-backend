package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	realtimegatewayadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/realtimegateway/adapters"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
)

var (
	DevelopmentRouter         *gin.Engine
	DevelopmentAPIRouterGroup *gin.RouterGroup
)

func ConfigureAPIRoutes(
	coreClient *coreadapters.CoreClient,
	realtimeGatewayClient *realtimegatewayadapters.RealtimeGatewayClient,
	allowedDomains []string,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {

	DevelopmentAPIRouterGroup = DevelopmentRouter.Group("/" + gatewaycontract.APIDevelopmentBaseURL) // use in development mode
	DevelopmentAPIRouterGroup.Use(
		middlewares.SanitizeXForwardedForMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.DomainWhiteListMiddleware(allowedDomains),
	)
	DevelopmentAPIRouterGroup.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(200) })
	fmt.Println("API router group path:", DevelopmentAPIRouterGroup.BasePath())

	configureDevelopmentAuthRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentUserRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentUserInfoRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureUserSettingRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentUserAccountRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentStationRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentRoutineRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentRoutineTagRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentRoutineTaskRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentRoutineTaskRecordRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentRootShelfRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentSubShelfRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentMaterialRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentBlockPackRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentBlockRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentRealtimeRoutes(DevelopmentAPIRouterGroup, coreClient, realtimeGatewayClient, accessTokenCookieHandler, refreshTokenCookieHandler)
	configureDevelopmentGraphQLRoutes(DevelopmentAPIRouterGroup, coreClient, accessTokenCookieHandler, refreshTokenCookieHandler)

	// test
	configureStaticRoutes(DevelopmentAPIRouterGroup)
}
