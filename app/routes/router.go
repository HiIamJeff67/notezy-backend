package routes

import (
	"fmt"
	"notezy-backend/app/controllers"
	constants "notezy-backend/global/constants"

	"github.com/gin-gonic/gin"
)

var (
	Router      *gin.Engine
	RouterGroup *gin.RouterGroup
)

func ConfigureRoutes() {
	RouterGroup = Router.Group("/api/" + constants.DevelopmentVersion) // use in development mode
	fmt.Println("Router group path:", RouterGroup.BasePath())

	configureUserRoutes()
}

func configureUserRoutes() {
	userRoutes := RouterGroup.Group("/user")
	{
		// userRoutes.GET("/me", func(ctx *gin.Context) {
		// 	ctx.JSON(200, gin.H{"message": "Hello, user!"})
		// })

		userRoutes.POST("/register", controllers.Register)
		userRoutes.POST("/login", controllers.Login)

		userRoutes.GET("/all", controllers.FindAllUsers)
	}
}
