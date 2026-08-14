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

func ConfigureAPIRoutes(
	coreClient *coreadapters.CoreAdapter,
	allowedDomains []string,
	rateLimiters RateLimiters,
) {
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
	configureDevelopmentStationRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
	configureDevelopmentRoutineRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
	configureDevelopmentRoutineTagRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
	configureDevelopmentRoutineTaskRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
	configureDevelopmentRootShelfRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
	configureDevelopmentSubShelfRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
	configureDevelopmentMaterialRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
	configureDevelopmentBlockPackRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
	configureDevelopmentBlockRoutes(DevelopmentAPIRouterGroup, coreClient, rateLimiters)
}
