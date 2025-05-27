package routes

import (
	"fmt"
	"go-gorm-api/app/controllers"
	"go-gorm-api/global"

	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
	RouterGroup *gin.RouterGroup
)

func ConfigureRoutes() {
    RouterGroup = Router.Group("/api/" + global.VERSION)
    fmt.Println("Router group path:", RouterGroup.BasePath())

    configureUserRoutes()
}

func configureUserRoutes() {
	userRoutes := RouterGroup.Group("/user")
	{
		userRoutes.GET("/me", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "Hello, user!"})
		})

		userRoutes.GET("/all", controllers.GetAllUsers)

		userRoutes.POST("/register", controllers.Register)
	}
}