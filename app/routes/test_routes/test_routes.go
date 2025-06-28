package testroutes

import (
	"fmt"
	"notezy-backend/shared/constants"

	"github.com/gin-gonic/gin"
)

var (
	TestRouter      *gin.Engine
	TestRouterGroup *gin.RouterGroup
)

func ConfigureTestRoutes() {
	TestRouterGroup = TestRouter.Group(constants.TestBaseURL)
	fmt.Println("Router group path:", TestRouterGroup.BasePath())

	ConfigureTestAuthRoutes(nil)
}
