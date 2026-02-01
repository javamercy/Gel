package cli

import (
	"github.com/spf13/cobra"
)

var (
	commitTreeMessageFlag string
	commitTreeParentsFlag []string
)
var commitTreeCmd = &cobra.Command{
	Use:   "commit-tree <tree-hash>",
	Short: "Create a new commit object from a tree object",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commitHash, err := commitTreeService.CommitTree(args[0], commitTreeMessageFlag, commitTreeParentsFlag)
		if err != nil {
			return err
		}

		cmd.Println(commitHash)
		return nil
	},
}

func init() {
	commitTreeCmd.Flags().StringVarP(&commitTreeMessageFlag, "message", "m", "", "Commit message")
	commitTreeCmd.Flags().StringSliceVarP(&commitTreeParentsFlag, "parent", "p", nil, "Parent commit(s)")
	rootCmd.AddCommand(commitTreeCmd)
}
