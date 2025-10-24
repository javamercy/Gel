package cmd

import (
	"Gel/src/gel/application/services"

	"github.com/spf13/cobra"
)

var updateIndexCmd = &cobra.Command{
	Use:     "update-index",
	Short:   "Update the index with the current state of the working directory",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		add, _ := cmd.Flags().GetBool("add")
		remove, _ := cmd.Flags().GetBool("remove")

		request := services.UpdateIndexRequest{
			Paths:  args,
			Add:    add,
			Remove: remove,
		}

		err := container.UpdateIndexService.UpdateIndex(request)
		if err != nil {
			cmd.PrintErrln("Error updating index:", err)
			return
		}
	},
}

func init() {
	updateIndexCmd.Flags().BoolP("add", "a", false, "Add files to the index")
	updateIndexCmd.Flags().BoolP("remove", "r", false, "Remove files from the index")
	rootCmd.AddCommand(updateIndexCmd)
}
