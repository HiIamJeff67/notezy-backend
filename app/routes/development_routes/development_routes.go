package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	constants "notezy-backend/shared/constants"
)

var (
	DevelopmentRouter      *gin.Engine
	DevelopmentRouterGroup *gin.RouterGroup
)

func ConfigureDevelopmentRoutes() {
	DevelopmentRouterGroup = DevelopmentRouter.Group(constants.DevelopmentBaseURL) // use in development mode
	fmt.Println("Router group path:", DevelopmentRouterGroup.BasePath())

	configureDevelopmentAuthRoutes()
	configureDevelopmentUserRoutes()
	configureDevelopmentUserInfoRoutes()
	configureDevelopmentUserAccountRoutes()
	configureDevelopmentGraphQLRoutes()
}
