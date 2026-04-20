package cli

import (
	"Gel/internal/staging"

	"github.com/spf13/cobra"
)

var (
	removeCachedFlag    bool
	removeDryRunFlag    bool
	removeRecursiveFlag bool
	removeForceFlag     bool
)

// removeCmd removes tracked files from the index and optionally from the working tree.
var removeCmd = &cobra.Command{
	Use:   "rm <pathspec>...",
	Short: "Remove a file or directory",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := removeService.Remove(
			args, staging.RemoveOptions{
				Cached:    removeCachedFlag,
				DryRun:    removeDryRunFlag,
				Recursive: removeRecursiveFlag,
				Force:     removeForceFlag,
			},
		)
		if err != nil {
			return err
		}

		for _, path := range result.Removed {
			cmd.Printf("rm '%s'\n", path)
		}
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVar(
		&removeCachedFlag, "cached", false,
		"Remove from cache only",
	)
	removeCmd.Flags().BoolVarP(
		&removeDryRunFlag, "dry-run", "n", false,
		"Show what would be removed without actually removing",
	)
	removeCmd.Flags().BoolVarP(
		&removeRecursiveFlag, "recursive", "r", false,
		"Remove directories and their contents recursively",
	)
	removeCmd.Flags().BoolVarP(
		&removeForceFlag, "force", "f", false,
		"Ignore nonexistent files and arguments, never prompt",
	)
	rootCmd.AddCommand(removeCmd)
}
