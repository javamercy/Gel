package cmd

import (
	"os"

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
			cmd.PrintErrln("ls-tree command requires exactly one argument: the tree hash")
			os.Exit(1)
		}
		hash := args[0]

		output, err := lsTreeService.LsTree(hash, recursiveFlag, showTreesFlag)
		if err != nil {
			return err
		}
		cmd.Println(output)
		return nil
	},
}

func init() {
	lsTreeCmd.SilenceUsage = true
	lsTreeCmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", false, "Recursively list subtrees")
	lsTreeCmd.Flags().BoolVarP(&showTreesFlag, "show-trees", "t", false, "Show tree objects in the listing")
	rootCmd.AddCommand(lsTreeCmd)
}
