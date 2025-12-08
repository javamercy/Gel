package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add <pathspec>...",
	Short:   "Add file contents to the index",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Error: no paths specified")
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")

		addedPaths, err := addService.Add(args, dryRun)
		if err != nil {
			cmd.PrintErrln("Error adding files:", err)
			return
		}

		if verbose || dryRun {
			for _, path := range addedPaths {
				if dryRun {
					cmd.Printf("add '%s'\n", path)
				} else {
					cmd.Println(path)
				}
			}
		}
	},
}

func init() {
	addCmd.Flags().BoolP("all", "A", false, "Add changes from all tracked and untracked files")
	addCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done, without making any changes")
	addCmd.Flags().BoolP("verbose", "v", false, "Show verbose output")
	rootCmd.AddCommand(addCmd)
}
