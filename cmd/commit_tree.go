package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	messageFlag string
)
var commitTreeCmd = &cobra.Command{
	Use:   "commit-tree <tree-hash>",
	Short: "Create a new commit object from a tree object",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("tree hash required")
		}

		treeHash := args[0]
		commitHash, err := commitTreeService.CommitTree(treeHash, messageFlag)
		if err != nil {
			return err
		}

		cmd.Println(commitHash)
		return nil
	},
}

func init() {
	commitTreeCmd.Flags().StringVarP(&messageFlag, "message", "m", "", "Commit message")
	rootCmd.AddCommand(commitTreeCmd)
}
