package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var readTreeCmd = &cobra.Command{
	Use:     "read-tree",
	Short:   "Read tree objects into the index",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			cmd.PrintErrln("Error: No tree hashes provided")
			os.Exit(1)
		}
		gelError := container.ReadTreeService.ReadTree(args[0])
		if gelError != nil {
			cmd.PrintErrln(gelError.Message)
			os.Exit(gelError.GetExitCode())
		}

	},
}

func init() {
	rootCmd.AddCommand(readTreeCmd)
}
