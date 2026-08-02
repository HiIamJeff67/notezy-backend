package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
)

var (
	DevelopmentRouter         *gin.Engine
	DevelopmentAPIRouterGroup *gin.RouterGroup
)

func ConfigureAPIRoutes() {
	coreClient := coreadapters.NewConfiguredCoreClient()

	DevelopmentAPIRouterGroup = DevelopmentRouter.Group("/" + constants.APIDevelopmentBaseURL) // use in development mode
	DevelopmentAPIRouterGroup.Use(
		middlewares.SanitizeXForwardedForMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.DomainWhiteListMiddleware(),
	)
	DevelopmentAPIRouterGroup.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(200) })
	fmt.Println("API router group path:", DevelopmentAPIRouterGroup.BasePath())

	configureDevelopmentAuthRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentUserRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentUserInfoRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureUserSettingRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentUserAccountRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentStationRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentRoutineRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentRoutineTagRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentRoutineTaskRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentRoutineTaskRecordRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentRootShelfRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentSubShelfRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentMaterialRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentBlockPackRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentBlockRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentRealtimeRoutes(DevelopmentAPIRouterGroup, coreClient)
	configureDevelopmentGraphQLRoutes(DevelopmentAPIRouterGroup, coreClient)

	// test
	configureStaticRoutes(DevelopmentAPIRouterGroup)
}
