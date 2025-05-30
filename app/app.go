package app

import (
	"fmt"
	models "notezy-backend/app/models"
	"notezy-backend/app/routes"
	"notezy-backend/global"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StartApplication() {
	models.NotezyDB = models.ConnectToDatabase(global.PostgresDatabaseConfig)

	routes.Router = gin.Default()
	routes.ConfigureRoutes()

	addr := global.GinAddr

	err := endless.ListenAndServe(addr, routes.Router)
	if err != nil {
		fmt.Println("Failed to connect to the server")
	}

	models.DisconnectToDatabase(models.NotezyDB)
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
