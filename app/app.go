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
	ReloadRedisLibraries()

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

func ReloadRedisLibraries() {
	if exception := caches.FlushCacheLibraries(); exception != nil {
		exception.Log()
	}
	if exception := caches.LoadRateLimitRecordCacheLibraries(); exception != nil {
		exception.Log()
	}
	if exception := caches.LoadUserQuotaCacheLibraries(); exception != nil {
		exception.Log()
	}
	// reload other more redis libraries here...
}
