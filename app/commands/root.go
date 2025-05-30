package commands

import (
	"notezy-backend/app"
	logs "notezy-backend/app/logs"
	models "notezy-backend/app/models"
	"notezy-backend/global"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "app",
	Short: "This is the root command.",
	Long:  "This is a longer description of the root command.",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Info("Welcome to the CLI.")
		app.StartApplication()
	},
}

var migrateDatabaseCommand = &cobra.Command{
	Use:   "migrateDB",
	Short: "Create or update database schema.",
	Long:  "Use models paclage to create or update database table schema.",
	Run: func(cmd *cobra.Command, args []string) {

		logs.Info("All avaliable databases:")
		for key, value := range models.DatabaseNameToInstance {
			logs.FInfo("database name: %v, instance: %v", key, value)
		}

		databaseNameStr, _ := cmd.Flags().GetString("database")
		if databaseNameStr == "" {
			logs.FError("The --database flag must be specified")
			return
		}

		db, ok := models.DatabaseNameToInstance[databaseNameStr]
		if !ok {
			logs.FError("The database instance with the name of %s is not exist", databaseNameStr)
			return
		}

		logs.Info("Start the process of migrating database schema.")
		app.MigrateDatabaseSchema(db)
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

		validTableName, isValidTableName := global.ConvertToValidTableName(tableNameStr)
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
		app.TrancateDatabaseTable(validTableName, db)
	},
}

// the Execute() function is the start point of cobra
func Execute() {
	/* register the migrate database command and its flags */
	migrateDatabaseCommand.Flags().String("database", "", "The name of the database to migrate")
	rootCommand.AddCommand(migrateDatabaseCommand)

	/* register the truncate database table command and its flags */
	truncateDatabaseCommand.Flags().String("database", "", "The name of the database to truncate the table inside it")
	truncateDatabaseCommand.Flags().String("table", "", "The name of the table to truncate")
	rootCommand.AddCommand(truncateDatabaseCommand)

	if err := rootCommand.Execute(); err != nil {
		logs.FError("Failed to init the CLI: %s", err)
		os.Exit(1)
	}
}
