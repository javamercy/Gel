package cmd

import "github.com/spf13/cobra"

var (
	messageFlag string
)
var commitTreeCmd = &cobra.Command{
	Use:     "commit-tree <tree-hash>",
	Short:   "Create a new commit object from a tree object",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Error: no tree hash specified")
			return
		}

		treeHash := args[0]

		commitHash, err := commitTreeService.CommitTree(treeHash, messageFlag)
		if err != nil {
			cmd.PrintErrln("Error creating commit:", err)
			return
		}
		cmd.Println(commitHash)
	},
}

func init() {
	commitTreeCmd.Flags().StringVarP(&messageFlag, "message", "m", "", "Commit message")
	rootCmd.AddCommand(commitTreeCmd)
}
