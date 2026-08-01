package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	configs "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/core/data"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

var viewAllAvailableDatabasesCommand = &cobra.Command{
	Use:   "viewDatabases",
	Short: "View all the available databases.",
	Long:  "Use some map to storing and printing the available databases in the project.",
	Run: func(cmd *cobra.Command, args []string) {
		logs.NotezyLogger.Info(context.Background(), "All available databases:")
		for key, value := range data.DatabaseNameToInstance {
			logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("database name: %v, instance: %v", key, value))
		}
	},
}

var viewAllDatabaseEnumsCommand = &cobra.Command{
	Use:   "viewAllEnums",
	Short: "View all the nums of the database.",
	Long:  "Use a simple select sql command to get all the enums of the database",
	Run: func(cmd *cobra.Command, args []string) {
		db := data.ConnectToDatabase(configs.PostgresDatabaseConfig)
		defer data.DisconnectToDatabase(db)

		if !data.ViewAllDatabaseEnums(db) {
			return
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
			logs.NotezyLogger.Error(context.Background(), nil, "The --database flag must be specified")
			return
		}

		tableNameStr, errorOfTableFlag := cmd.Flags().GetString("table")
		if errorOfTableFlag != nil {
			logs.NotezyLogger.Error(context.Background(), nil, "The --table flag must be specified")
			return
		}

		tableName, isTableName := types.ConvertToTableName(tableNameStr)
		if !isTableName {
			logs.NotezyLogger.Error(context.Background(), nil, fmt.Sprintf("The table name of %s is not in the database %s", tableNameStr, databaseNameStr))
			return
		}

		db, ok := data.DatabaseNameToInstance[tableNameStr]
		if !ok {
			logs.NotezyLogger.Error(context.Background(), nil, "The database instance is not exist")
			return
		}

		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Start the process of truncating database table: %s", tableNameStr))
		db = data.ConnectToDatabase(data.DatabaseInstanceToConfig[db])
		defer data.DisconnectToDatabase(db)

		data.TruncateTablesInDatabase(tableName, db)
	},
}

var migrateDatabaseCommand = &cobra.Command{
	Use:   "migrateDB",
	Short: "Migrate enums, tables, and some triggers to the database.",
	Long:  "Use some migration SQLs to migrate required enums, tables, and some triggers to the database.",
	Run: func(cmd *cobra.Command, args []string) {
		db := data.ConnectToDatabase(configs.PostgresDatabaseConfig)
		defer data.DisconnectToDatabase(db)

		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Start the process of migrating database schema to %v", configs.PostgresDatabaseConfig.DBName))

		if !data.MigrateEnumsToDatabase(db) {
			return
		}
		if !data.MigrateTablesToDatabase(db) {
			return
		}
		if !data.MigrateTriggersToDatabase(db) {
			return
		}
		if !data.MigrateConstraintsToDatabase(db) {
			return
		}
	},
}

var seedDatabaseCommand = &cobra.Command{
	Use:   "seedDB",
	Short: "Seed some default data for management or main business logic.",
	Long:  "Use some seeding default data SQLs to seed data for management or main business logic.",
	Run: func(cmd *cobra.Command, args []string) {
		db := data.ConnectToDatabase(configs.PostgresDatabaseConfig)
		defer data.DisconnectToDatabase(db)

		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Start the process of seeding database default data to %v", configs.PostgresDatabaseConfig.DBName))

		if !data.SeedDefaultDataToDatabase(db) {
			return
		}
	},
}

/* ============================== Prepare Flags Helper Function ============================== */

func PrepareDatabaseCommandsFlags() {
	/* register the flags of truncating database table command */
	truncateDatabaseCommand.Flags().String("database", "", "The name of the database to truncate the table inside it")
	truncateDatabaseCommand.Flags().String("table", "", "The name of the table to truncate")
}
