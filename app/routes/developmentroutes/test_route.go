package developmentroutes

import (
	"notezy-backend/app/monitor/logs"

	"github.com/gin-gonic/gin"
)

func configureDevelopmentTestRoutes() {
	testRoutes := DevelopmentRouterGroup.Group("/test")
	{
		testRoutes.POST(
			"/webhook",
			func(ctx *gin.Context) {
				body, err := ctx.GetRawData()
				if err != nil {
					logs.Error("failed to read request body: ", err)
					ctx.JSON(400, gin.H{"error": "Invalid request"})
					return
				}

				logs.Info(string(body))
				ctx.JSON(200, gin.H{"message": "webhook received!"})
			},
		)
	}
}
