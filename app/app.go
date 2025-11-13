package app

import (
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	caches "notezy-backend/app/caches"
	models "notezy-backend/app/models"
	developmentroutes "notezy-backend/app/routes/development_routes"
	util "notezy-backend/app/util"
)

func StartApplication() {
	models.NotezyDB = models.ConnectToDatabase(models.PostgresDatabaseConfig)
	caches.ConnectToAllRedis()
	ReLoadRedisFunctions()

	developmentroutes.DevelopmentRouter = gin.Default()
	developmentroutes.ConfigureDevelopmentRoutes()

	ginAddr := util.GetEnv("GIN_DOMAIN", "") + ":" + util.GetEnv("GIN_PORT", "7777")

	err := endless.ListenAndServe(ginAddr, developmentroutes.DevelopmentRouter)
	if err != nil {
		fmt.Println("Failed to connect to the server")
	}

	models.DisconnectToDatabase(models.NotezyDB)
	caches.DisconnectToAllRedis()
}

func ReLoadRedisFunctions() {
	if exception := caches.FlushRedisFunctionsLibraries(); exception != nil {
		exception.Log()
	}
	if exception := caches.ReloadRateLimitRecordRedisFunctionsLibraries(); exception != nil {
		exception.Log()
	}
}
