package commands

import (
	"os"

	"github.com/spf13/cobra"

	app "notezy-backend/app"
	logs "notezy-backend/app/logs"
	traces "notezy-backend/app/traces"
)

var rootCommand = &cobra.Command{
	Use:   "app",
	Short: "This is the root command.",
	Long:  "This is a longer description of the root command.",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Info(traces.GetTrace(0).FileLineString(), "Welcome to the CLI.")
		app.StartApplication()
	},
}

func AddCommands(rootCommand *cobra.Command, addedCommands []*cobra.Command) {
	for _, command := range addedCommands {
		rootCommand.AddCommand(command)
	}
}

// the Execute() function is the start point of cobra
func Execute() {
	// prepare the flags of database commands
	PrepareDatabaseCommandsFlags()
	// add the commands of database
	AddCommands(
		rootCommand,
		[]*cobra.Command{
			viewAllAvailableDatabasesCommand,
			truncateDatabaseCommand,
			viewAllDatabaseEnumsCommand,
			migrateDatabaseCommand,
			seedDatabaseCommand,
		},
	)

	if err := rootCommand.Execute(); err != nil {
		logs.FError(traces.GetTrace(0).FileLineString(), "Failed to init the CLI: %s", err)
		os.Exit(1)
	}
}
