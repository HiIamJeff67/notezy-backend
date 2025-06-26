package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	constants "notezy-backend/shared/constants"
)

var (
	Router      *gin.Engine
	RouterGroup *gin.RouterGroup
)

func ConfigureRoutes() {
	RouterGroup = Router.Group(constants.BaseURL) // use in development mode
	fmt.Println("Router group path:", RouterGroup.BasePath())

	configureAuthRoutes()
	configureUserRoutes()
	configureTestRoutes()
	configureUserInfoRoutes()
	configureUserAccountRoutes()
}
