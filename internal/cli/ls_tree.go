package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	recursiveFlag bool
	showTreesFlag bool
)
var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("ls-tree command requires tree hash")
		}

		hash := args[0]
		return lsTreeService.LsTree(cmd.OutOrStdout(), hash, recursiveFlag, showTreesFlag)

	},
}

func init() {
	lsTreeCmd.SilenceUsage = true
	lsTreeCmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", false, "Recursively list subtrees")
	lsTreeCmd.Flags().BoolVarP(&showTreesFlag, "show-trees", "t", false, "Show tree objects in the listing")
	rootCmd.AddCommand(lsTreeCmd)
}
