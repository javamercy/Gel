package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dryRunFlag bool
)
var addCmd = &cobra.Command{
	Use:   "add <pathspec>...",
	Short: "Add file contents to the index",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("at least one pathspec required")
		}

		return addService.Add(cmd.OutOrStdout(), args, dryRunFlag)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "n", false, "Dry run the add operation without making any changes")
	rootCmd.AddCommand(addCmd)
}
