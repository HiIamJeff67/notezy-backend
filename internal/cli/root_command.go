package main

import "github.com/spf13/cobra"

var rootCommand = &cobra.Command{
	Use:   "notezy",
	Short: "Notezy development and verification commands.",
}

func Execute() error {
	return rootCommand.Execute()
}
