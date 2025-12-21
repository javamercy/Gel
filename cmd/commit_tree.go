package cmd

import "github.com/spf13/cobra"

var commitTreeCmd = &cobra.Command{
	Use:     "commit-tree <tree-hash>",
	Short:   "Create a new commit object from a tree object",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Error: no tree hash specified")
			return
		}
	},
}

func init() {
	commitTreeCmd.Flags().StringP("tree-hash", "t", "", "Tree hash")
	commitTreeCmd.Flags().StringP("message", "m", "", "Commit message")
	rootCmd.AddCommand(commitTreeCmd)
}
