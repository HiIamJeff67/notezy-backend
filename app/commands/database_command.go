package commands

import (
	"github.com/spf13/cobra"

	logs "notezy-backend/app/logs"
	models "notezy-backend/app/models"
	types "notezy-backend/shared/types"
)

var viewAllAvailableDatabasesCommand = &cobra.Command{
	Use:   "viewDatabases",
	Short: "View all the available databases.",
	Long:  "Use some map to storing and printing the available databases in the project.",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Info("All available databases:")
		for key, value := range models.DatabaseNameToInstance {
			logs.FInfo("database name: %v, instance: %v", key, value)
		}
	},
}

var migrateDatabaseCommand = &cobra.Command{
	Use:   "migrateDB",
	Short: "Create or update database schema.",
	Long:  "Use models package to create or update database table schema.",
	Run: func(cmd *cobra.Command, args []string) {
		db := models.ConnectToDatabase(models.PostgresDatabaseConfig)
		logs.FInfo("Start the process of migrating database schema to %v.", models.PostgresDatabaseConfig.DBName)
		if !models.MigrateEnumsToDatabase(db) {
			return
		}
		if !models.MigrateTablesToDatabase(db) {
			return
		}
		if !models.MigrateTriggersToDatabase(db) {
			return
		}
		models.DisconnectToDatabase(db)
	},
}

var truncateDatabaseCommand = &cobra.Command{
	Use:   "truncate",
	Short: "Truncate an existing table",
	Long:  "Truncate the database table with the given table name",
	Run: func(cmd *cobra.Command, args []string) {
		databaseNameStr, errorOfDatabaseFlag := cmd.Flags().GetString("database")
		if errorOfDatabaseFlag != nil {
			logs.FError("The --database flag must be specified")
			return
		}

		tableNameStr, errorOfTableFlag := cmd.Flags().GetString("table")
		if errorOfTableFlag != nil {
			logs.FError("The --table flag must be specified")
			return
		}

		validTableName, isValidTableName := types.ConvertToValidTableName(tableNameStr)
		if !isValidTableName {
			logs.FError("The table name of %s is not in the database %s", tableNameStr, databaseNameStr)
			return
		}

		db, ok := models.DatabaseNameToInstance[tableNameStr]
		if !ok {
			logs.FError("The database instance is not exist")
			return
		}

		logs.FInfo("Start the process of truncating database table: %s.", tableNameStr)
		models.NotezyDB = models.ConnectToDatabase(models.DatabaseInstanceToConfig[db])
		models.TruncateTablesInDatabase(validTableName, db)
		models.DisconnectToDatabase(models.NotezyDB)
	},
}

/* ============================== Parepare Flags Helper Function ============================== */

func PrepareDatabaseCommandsFlags() {
	/* register the flags of truncating database table command */
	truncateDatabaseCommand.Flags().String("database", "", "The name of the database to truncate the table inside it")
	truncateDatabaseCommand.Flags().String("table", "", "The name of the table to truncate")
}
