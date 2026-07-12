package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

var (
	DevelopmentRouter         *gin.Engine
	DevelopmentAPIRouterGroup *gin.RouterGroup
)

func ConfigureAPIRoutes() {
	DevelopmentAPIRouterGroup = DevelopmentRouter.Group("/" + constants.APIDevelopmentBaseURL) // use in development mode
	DevelopmentAPIRouterGroup.Use(
		middlewares.SanitizeXForwardedForMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.DomainWhiteListMiddleware(),
	)
	DevelopmentAPIRouterGroup.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(200) })
	fmt.Println("API router group path:", DevelopmentAPIRouterGroup.BasePath())

	configureDevelopmentAuthRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentUserRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentUserInfoRoutes(DevelopmentAPIRouterGroup)
	configureUserSettingRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentUserAccountRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentStationRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentRoutineRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentRoutineTagRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentRoutineTaskRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentRoutineTaskRecordRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentRootShelfRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentSubShelfRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentMaterialRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentBlockPackRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentRealtimeAPIRoutes(DevelopmentAPIRouterGroup)
	configureDevelopmentGraphQLRoutes(DevelopmentAPIRouterGroup)

	// test
	configureStaticRoutes(DevelopmentAPIRouterGroup)
	configureStorageRoutes(DevelopmentAPIRouterGroup)
}
