package app

import (
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	caches "notezy-backend/app/caches"
	models "notezy-backend/app/models"
	developmentroutes "notezy-backend/app/routes/development_routes"
	shared "notezy-backend/shared"
)

func StartApplication() {
	models.NotezyDB = models.ConnectToDatabase(shared.PostgresDatabaseConfig)
	caches.ConnectToAllRedis()

	developmentroutes.DevelopmentRouter = gin.Default()
	developmentroutes.ConfigureDevelopmentRoutes()

	addr := shared.GinAddr

	err := endless.ListenAndServe(addr, developmentroutes.DevelopmentRouter)
	if err != nil {
		fmt.Println("Failed to connect to the server")
	}

	models.DisconnectToDatabase(models.NotezyDB)
	caches.DisconnectToAllRedis()
}

func MigrateDatabaseSchema(db *gorm.DB) {
	localDB := models.ConnectToDatabase(models.DatabaseInstanceToConfig[db])
	models.MigrateToDatabase(localDB)
	models.DisconnectToDatabase(localDB)
}

func TrancateDatabaseTable(tableName shared.ValidTableName, db *gorm.DB) {
	models.NotezyDB = models.ConnectToDatabase(models.DatabaseInstanceToConfig[db])
	models.TruncateTablesInDatabase(tableName, db)
	models.DisconnectToDatabase(models.NotezyDB)
}
