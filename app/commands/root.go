package commands

import (
	"os"

	"github.com/spf13/cobra"

	app "notezy-backend/app"
	logs "notezy-backend/app/logs"
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

// the Execute() function is the start point of cobra
func Execute() {
	/* register the view all available databases command */
	rootCommand.AddCommand(viewAllAvailableDatabasesCommand)

	/* register the migrate database command */
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
