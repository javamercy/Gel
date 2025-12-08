package cmd

import (
	"github.com/spf13/cobra"
)

var readTreeCmd = &cobra.Command{
	Use:     "read-tree <tree-hash>",
	Short:   "Read tree objects into the index",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Error: tree hash required")
			return
		}

		hash := args[0]
		err := readTreeService.ReadTree(hash)
		if err != nil {
			cmd.PrintErrln("Error reading tree:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(readTreeCmd)
}
