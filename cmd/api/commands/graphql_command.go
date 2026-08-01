package commands

import "github.com/spf13/cobra"

var writeGraphQLEnumMappingValuesToConfig = &cobra.Command{
	Use:   "writeGQLEnumMappingValues",
	Short: "Write gql enum mapping values to config file.",
	Long:  "Use truth enum values defined in golang and database to write them as enum mapping values to the graphql config file.",
	Run: func(cmd *cobra.Command, args []string) {
		// may be implemented in the future
	},
}

/* ============================== Prepare Flags Helper Function ============================== */

func PrepareGraphQLCommandsFlags() {

}
