package testroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

var (
	TestRouter      *gin.Engine
	TestRouterGroup *gin.RouterGroup
)

func ConfigureTestRoutes(db *gorm.DB) {
	TestRouterGroup = TestRouter.Group("/" + constants.TestBaseURL)
	fmt.Println("Router group path:", TestRouterGroup.BasePath())

	ConfigureTestAuthRoutes(db, TestRouterGroup)
}
