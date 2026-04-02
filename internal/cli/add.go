package cli

import (
	"Gel/internal/staging"

	"github.com/spf13/cobra"
)

var (
	addDryRunFlag  bool
	addVerboseFlag bool
)

// addCmd stages file content into the index using pathspec semantics.
var addCmd = &cobra.Command{
	Use:   "add <pathspec>...",
	Short: "Add file contents to the index",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addResult := addService.Add(
			args,
			staging.AddOptions{
				Verbose: addVerboseFlag,
				DryRun:  addDryRunFlag,
			},
		)
		if addResult.Error != nil {
			return addResult.Error
		}
		for _, file := range addResult.Added {
			cmd.Printf("A %s\n", file)
		}
		for _, file := range addResult.Removed {
			cmd.Printf("D %s\n", file)
		}
		return nil
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
