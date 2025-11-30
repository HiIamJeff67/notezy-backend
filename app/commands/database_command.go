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

		tableName, isTableName := types.ConvertToTableName(tableNameStr)
		if !isTableName {
			logs.FError("The table name of %s is not in the database %s", tableNameStr, databaseNameStr)
			return
		}

		db, ok := models.DatabaseNameToInstance[tableNameStr]
		if !ok {
			logs.FError("The database instance is not exist")
			return
		}

		logs.FInfo("Start the process of truncating database table: %s", tableNameStr)
		models.NotezyDB = models.ConnectToDatabase(models.DatabaseInstanceToConfig[db])
		models.TruncateTablesInDatabase(tableName, db)

		models.DisconnectToDatabase(models.NotezyDB)
	},
}

var migrateDatabaseCommand = &cobra.Command{
	Use:   "migrateDB",
	Short: "Migrate enums, tables, and some triggers to the database.",
	Long:  "Use some migration SQLs to migrate required enums, tables, and some triggers to the database.",
	Run: func(cmd *cobra.Command, args []string) {
		db := models.ConnectToDatabase(models.PostgresDatabaseConfig)

		logs.FInfo("Start the process of migrating database schema to %v", models.PostgresDatabaseConfig.DBName)

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

var seedDatabaseCommand = &cobra.Command{
	Use:   "seedDB",
	Short: "Seed some default data for management or main business logic.",
	Long:  "Use some seeding default data SQLs to seed data for management or main business logic.",
	Run: func(cmd *cobra.Command, args []string) {
		db := models.ConnectToDatabase(models.PostgresDatabaseConfig)

		logs.FInfo("Start the process of seeding database default data to %v", models.PostgresDatabaseConfig.DBName)

		if !models.SeedDefaultDataToDatabase(db) {
			return
		}

		models.DisconnectToDatabase(db)
	},
}

/* ============================== Parepare Flags Helper Function ============================== */

func PrepareDatabaseCommandsFlags() {
	/* register the flags of truncating database table command */
	truncateDatabaseCommand.Flags().String("database", "", "The name of the database to truncate the table inside it")
	truncateDatabaseCommand.Flags().String("table", "", "The name of the table to truncate")
}
