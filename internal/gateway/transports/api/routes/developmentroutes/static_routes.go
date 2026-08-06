package developmentroutes

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"

	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
)

func configureStaticRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	staticGroup := router.Group("/static")
	{
		globalImagesGroup := staticGroup.Group("/globalImages")
		globalImagesGroup.Use(
			middlewares.UnauthorizedRateLimitMiddleware(),
		)
		{
			// configure avatars
			globalImagesGroup.GET("/avatars/:id", func(ctx *gin.Context) {
				ctx.Header("Cross-Origin-Resource-Policy", "cross-origin")
				avatarId := ctx.Param("id")
				filePath := fmt.Sprintf("./global/images/avatars/userAvatar%s.png", avatarId)

				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					filePath = "./global/images/avatars/userAvatar1.png"
				}
				logs.NotezyLogger.Info(ctx.Request.Context(), "download file")

				ctx.File(filePath)
			})

			// configure brand icon here in the future
		}
	}
}
