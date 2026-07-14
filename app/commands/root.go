package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
)

var rootCommand = &cobra.Command{
	Use:   "app",
	Short: "This is the root command.",
	Long:  "This is a longer description of the root command.",
}

func AddCommands(rootCommand *cobra.Command, addedCommands []*cobra.Command) {
	for _, command := range addedCommands {
		rootCommand.AddCommand(command)
	}
}

// the Execute() function is the start point of cobra
func Execute() {
	logs.NotezyLogger = logs.NewCommandLineInterfaceLogger()

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
		logs.NotezyLogger.Error(context.Background(), nil, fmt.Sprintf("Failed to init the CLI: %s", err))
		os.Exit(1)
	}
}
