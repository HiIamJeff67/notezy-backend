package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/apigateway/ratelimit"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/core/adapters"
)

var (
	DevelopmentRouter         *gin.Engine
	DevelopmentAPIRouterGroup *gin.RouterGroup
)

type RateLimiters struct {
	Unauthorized *ratelimit.HybridRateLimiter
}

type APIRouteDependencies struct {
	CoreClient     *coreadapters.CoreAdapter
	AllowedDomains []string
	RateLimiters   RateLimiters
}

func NewRouter(deps APIRouteDependencies) *gin.Engine {
	DevelopmentRouter = gin.Default()
	coreClient, allowedDomains, rateLimiters := deps.CoreClient, deps.AllowedDomains, deps.RateLimiters
	DevelopmentAPIRouterGroup = DevelopmentRouter.Group("/" + gatewaycontract.APIDevelopmentBaseURL) // use in development mode
	DevelopmentAPIRouterGroup.Use(
		middlewares.SanitizeXForwardedForMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.DomainWhiteListMiddleware(allowedDomains),
	)
	DevelopmentAPIRouterGroup.Use(middlewares.KeyMiddleware())
	DevelopmentAPIRouterGroup.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(200) })
	fmt.Println("API router group path:", DevelopmentAPIRouterGroup.BasePath())

	// APIGateway deliberately exposes only the first-party resource domains
	// that are stable enough for external integrations.
	configureDevelopmentStationRoutes(DevelopmentAPIRouterGroup, StationRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})
	configureDevelopmentRoutineRoutes(DevelopmentAPIRouterGroup, RoutineRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})
	configureDevelopmentRoutineTagRoutes(DevelopmentAPIRouterGroup, RoutineTagRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})
	configureDevelopmentRoutineTaskRoutes(DevelopmentAPIRouterGroup, RoutineTaskRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})
	configureDevelopmentRootShelfRoutes(DevelopmentAPIRouterGroup, RootShelfRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})
	configureDevelopmentSubShelfRoutes(DevelopmentAPIRouterGroup, SubShelfRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})
	configureDevelopmentMaterialRoutes(DevelopmentAPIRouterGroup, MaterialRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})
	configureDevelopmentBlockPackRoutes(DevelopmentAPIRouterGroup, BlockPackRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})
	configureDevelopmentBlockRoutes(DevelopmentAPIRouterGroup, BlockRouteDependencies{CoreClient: coreClient, RateLimiters: rateLimiters})

	return DevelopmentRouter
}
