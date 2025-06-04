package routes

import (
	"fmt"
	constants "notezy-backend/global/constants"

	"github.com/gin-gonic/gin"
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
}
