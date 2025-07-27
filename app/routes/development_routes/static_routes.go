package developmentroutes

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	logs "notezy-backend/app/logs"
)

func configureStaticRoutes() {
	staticGroup := DevelopmentRouterGroup.Group("/static")
	{
		globalImagesGroup := staticGroup.Group("/globalImages")
		{
			// configure avatars
			globalImagesGroup.GET("/avatars/:id", func(ctx *gin.Context) {
				avatarId := ctx.Param("id")
				filePath := fmt.Sprintf("./global/images/avatars/userAvatar%s.png", avatarId)

				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					filePath = "./global/images/avatars/userAvatar1.png"
				}
				logs.FInfo("download file")

				ctx.File(filePath)
			})

			// configure brand icon here in the future
		}
	}
}
