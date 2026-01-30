package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	addFlag    bool
	removeFlag bool
)
var updateIndexCmd = &cobra.Command{
	Use:   "update-index <file>...",
	Short: "Update the index with the current state of the working directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no paths specified")
		}

		if !addFlag && !removeFlag {
			return fmt.Errorf("must specify either --add or --remove")
		}

		return updateIndexService.UpdateIndex(args, addFlag, removeFlag)
	},
}

func init() {
	updateIndexCmd.Flags().BoolVarP(&addFlag, "add", "a", false, "Add specified files to the index")
	updateIndexCmd.Flags().BoolVarP(&removeFlag, "remove", "r", false, "Remove specified files from the index")
	rootCmd.AddCommand(updateIndexCmd)
}
