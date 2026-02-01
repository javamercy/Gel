package cli

import (
	"github.com/spf13/cobra"
)

var (
	lsTreeRecursiveFlag bool
	lsTreeShowTreesFlag bool
	lsTreeNameOnlyFlag  bool
)
var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := args[0]
		return lsTreeService.LsTree(
			cmd.OutOrStdout(), hash, lsTreeRecursiveFlag, lsTreeShowTreesFlag, lsTreeNameOnlyFlag,
		)
	},
}

func init() {
	lsTreeCmd.Flags().BoolVarP(&lsTreeRecursiveFlag, "recursive", "r", false, "Recursively list subtrees")
	lsTreeCmd.Flags().BoolVarP(&lsTreeShowTreesFlag, "show-trees", "t", false, "Show tree objects in the listing")
	lsTreeCmd.Flags().BoolVarP(&lsTreeNameOnlyFlag, "name-only", "n", false, "Show only names of the entries")
	rootCmd.AddCommand(lsTreeCmd)
}
