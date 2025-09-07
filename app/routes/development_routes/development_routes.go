package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	middlewares "notezy-backend/app/middlewares"
	constants "notezy-backend/shared/constants"
)

var (
	DevelopmentRouter      *gin.Engine
	DevelopmentRouterGroup *gin.RouterGroup
)

func ConfigureDevelopmentRoutes() {
	DevelopmentRouterGroup = DevelopmentRouter.Group("/" + constants.DevelopmentBaseURL) // use in development mode
	DevelopmentRouterGroup.Use(
		middlewares.CORSMiddleware(),
		middlewares.DomainWhitelistMiddleware(),
	)
	DevelopmentRouterGroup.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(200) })
	fmt.Println("Router group path:", DevelopmentRouterGroup.BasePath())

	configureStaticRoutes()
	configureDevelopmentAuthRoutes()
	configureDevelopmentUserRoutes()
	configureDevelopmentUserInfoRoutes()
	configureUserSettingRoutes()
	configureDevelopmentUserAccountRoutes()
	configureDevelopmentGraphQLRoutes()
	configureDevelopmentShelfRoutes()
	configureDevelopmentMaterialRoutes()

	// test
	configureStorageRoutes()
}
