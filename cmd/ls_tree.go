package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var lsTreeCmd = &cobra.Command{
	Use:     "ls-tree",
	Short:   "List the contents of a tree",
	PreRunE: requiresEnsureContextPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.PrintErrln("ls-tree command requires exactly one argument: the tree hash")
			os.Exit(1)
		}
		hash := args[0]

		recursive, _ := cmd.Flags().GetBool("recursive")
		showTrees, _ := cmd.Flags().GetBool("show-trees")

		output, err := lsTreeService.LsTree(hash, recursive, showTrees)
		if err != nil {
			return err
		}
		cmd.Println(output)
		return nil
	},
}

func init() {
	lsTreeCmd.SilenceUsage = true
	lsTreeCmd.Flags().BoolP("recursive", "r", false, "")
	lsTreeCmd.Flags().BoolP("show-trees", "t", false, "")
	rootCmd.AddCommand(lsTreeCmd)
}
