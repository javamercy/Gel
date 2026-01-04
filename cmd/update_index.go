package cmd

import (
	"github.com/spf13/cobra"
)

var (
	addFlag    bool
	removeFlag bool
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

		if !addFlag && !removeFlag {
			cmd.PrintErrln("Error: must specify either --add or --remove")
			return
		}

		err := updateIndexService.UpdateIndex(args, addFlag, removeFlag)
		if err != nil {
			cmd.PrintErrln("Error updating index:", err)
			return
		}
	},
}

func init() {
	updateIndexCmd.Flags().BoolVarP(&addFlag, "add", "a", false, "Add specified files to the index")
	updateIndexCmd.Flags().BoolVarP(&removeFlag, "remove", "r", false, "Remove specified files from the index")
	rootCmd.AddCommand(updateIndexCmd)
}
