package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:     "write-tree",
	Short:   "Write the current index as a tree object",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		hash, gelError := container.WriteTreeService.WriteTree()
		if gelError != nil {
			cmd.PrintErrln(gelError)
			os.Exit(gelError.GetExitCode())
		}
		cmd.Println(hash)
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
