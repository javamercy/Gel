package cmd

import (
	"github.com/spf13/cobra"
)

var updateIndexCmd = &cobra.Command{
	Use:     "update-index <file>...",
	Short:   "Update the index with the current state of the working directory",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Error: no paths specified")
			return
		}

		add, _ := cmd.Flags().GetBool("add")
		remove, _ := cmd.Flags().GetBool("remove")

		if !add && !remove {
			cmd.PrintErrln("Error: must specify either --add or --remove")
			return
		}

		err := updateIndexService.UpdateIndex(args, add, remove)
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
