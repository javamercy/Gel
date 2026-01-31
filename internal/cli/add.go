package cli

import (
	"github.com/spf13/cobra"
)

var (
	addDryRunFlag  bool
	addVerboseFlag bool
)
var addCmd = &cobra.Command{
	Use:   "add <pathspec>...",
	Short: "Add file contents to the index",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addService.Add(cmd.OutOrStdout(), args, addDryRunFlag, addVerboseFlag)
	},
}

func init() {
	addCmd.Flags().BoolVarP(
		&addDryRunFlag, "dry-run", "n", false,
		"Dry run the add operation without making any changes",
	)
	addCmd.Flags().BoolVarP(
		&addVerboseFlag, "verbose", "v", false,
		"Show verbose output of the add operation",
	)
	rootCmd.AddCommand(addCmd)
}
