package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	platformdatabase "github.com/HiIamJeff67/notezy-backend/shared/platform/database"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"

	data "github.com/HiIamJeff67/notezy-backend/internal/core/data/database"
)

var viewAllAvailableDatabasesCommand = &cobra.Command{
	Use:   "viewDatabases",
	Short: "View all the available databases.",
	Long:  "Use some map to storing and printing the available databases in the project.",
	Run: func(_ *cobra.Command, _ []string) {
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
	Run: func(_ *cobra.Command, _ []string) {
		config, err := platformdatabase.LoadConfig()
		if err != nil {
			panic(err)
		}
		db := data.ConnectToDatabase(config)
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
	Run: func(command *cobra.Command, _ []string) {
		databaseName, err := command.Flags().GetString("database")
		if err != nil {
			logs.NotezyLogger.Error(context.Background(), nil, "The --database flag must be specified")
			return
		}

		tableNameString, err := command.Flags().GetString("table")
		if err != nil {
			logs.NotezyLogger.Error(context.Background(), nil, "The --table flag must be specified")
			return
		}

		tableName, exists := data.ConvertToTableName(tableNameString)
		if !exists {
			logs.NotezyLogger.Error(context.Background(), nil, fmt.Sprintf("The table name of %s is not in the database %s", tableNameString, databaseName))
			return
		}

		db, exists := data.DatabaseNameToInstance[tableNameString]
		if !exists {
			logs.NotezyLogger.Error(context.Background(), nil, "The database instance is not exist")
			return
		}

		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Start the process of truncating database table: %s", tableNameString))
		db = data.ConnectToDatabase(data.DatabaseInstanceToConfig[db])
		defer data.DisconnectToDatabase(db)

		data.TruncateTablesInDatabase(tableName, db)
	},
}

var migrateDatabaseCommand = &cobra.Command{
	Use:   "migrateDB",
	Short: "Migrate enums, tables, and some triggers to the database.",
	Long:  "Use some migration SQLs to migrate required enums, tables, and some triggers to the database.",
	Run: func(_ *cobra.Command, _ []string) {
		config, err := platformdatabase.LoadConfig()
		if err != nil {
			panic(err)
		}
		db := data.ConnectToDatabase(config)
		defer data.DisconnectToDatabase(db)

		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Start the process of migrating database schema to %v", config.Name))

		if !data.MigrateEnumsToDatabase(db) ||
			!data.MigrateTablesToDatabase(db) ||
			!data.MigrateTriggersToDatabase(db) ||
			!data.MigrateConstraintsToDatabase(db) {
			return
		}
	},
}

var seedDatabaseCommand = &cobra.Command{
	Use:   "seedDB",
	Short: "Seed some default data for management or main business logic.",
	Long:  "Use some seeding default data SQLs to seed data for management or main business logic.",
	Run: func(_ *cobra.Command, _ []string) {
		config, err := platformdatabase.LoadConfig()
		if err != nil {
			panic(err)
		}
		db := data.ConnectToDatabase(config)
		defer data.DisconnectToDatabase(db)

		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Start the process of seeding database default data to %v", config.Name))

		if !data.SeedDefaultDataToDatabase(db) {
			return
		}
	},
}
