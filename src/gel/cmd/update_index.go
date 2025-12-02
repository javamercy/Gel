package cmd

import (
	"Gel/src/gel/application/dto"
	"os"

	"github.com/spf13/cobra"
)

var updateIndexCmd = &cobra.Command{
	Use:     "update-index",
	Short:   "Update the index with the current state of the working directory",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		add, _ := cmd.Flags().GetBool("add")
		remove, _ := cmd.Flags().GetBool("remove")

		request := dto.NewUpdateIndexRequest(args, add, remove)

		gelError := container.UpdateIndexService.UpdateIndex(request)
		if gelError != nil {
			cmd.PrintErrln(gelError)
			os.Exit(gelError.GetExitCode())
		}
	},
}

func init() {
	updateIndexCmd.Flags().BoolP("add", "a", false, "Add files to the index")
	updateIndexCmd.Flags().BoolP("remove", "r", false, "Remove files from the index")
	rootCmd.AddCommand(updateIndexCmd)
}
