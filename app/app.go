package app

import (
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	caches "notezy-backend/app/caches"
	models "notezy-backend/app/models"
	developmentroutes "notezy-backend/app/routes/development_routes"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

func StartApplication() {
	models.NotezyDB = models.ConnectToDatabase(models.PostgresDatabaseConfig)
	caches.ConnectToAllRedis()

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

func MigrateDatabaseSchemas(db *gorm.DB) {
	// execute the below migrations in sequence

	if !models.MigrateEnumsToDatabase(db) {
		return
	}
	if !models.MigrateTablesToDatabase(db) {
		return
	}
	if !models.MigrateTriggersToDatabase(db) {
		return
	}
}

func TrancateDatabaseTable(tableName types.ValidTableName, db *gorm.DB) {
	models.NotezyDB = models.ConnectToDatabase(models.DatabaseInstanceToConfig[db])
	models.TruncateTablesInDatabase(tableName, db)
	models.DisconnectToDatabase(models.NotezyDB)
}
