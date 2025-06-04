package app

import (
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	caches "notezy-backend/app/caches"
	models "notezy-backend/app/models"
	routes "notezy-backend/app/routes"
	global "notezy-backend/global"
)

func StartApplication() {
	models.NotezyDB = models.ConnectToDatabase(global.PostgresDatabaseConfig)
	caches.ConnectToAllRedis()

	routes.Router = gin.Default()
	routes.ConfigureRoutes()

	addr := global.GinAddr

	err := endless.ListenAndServe(addr, routes.Router)
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

func TrancateDatabaseTable(tableName global.ValidTableName, db *gorm.DB) {
	models.NotezyDB = models.ConnectToDatabase(models.DatabaseInstanceToConfig[db])
	models.TruncateTablesInDatabase(tableName, db)
	models.DisconnectToDatabase(models.NotezyDB)
}
