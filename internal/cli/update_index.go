package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	updateIndexAddFlag    bool
	updateIndexRemoveFlag bool
)
var updateIndexCmd = &cobra.Command{
	Use:   "update-index <file>...",
	Short: "Update the index with the current state of the working directory",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no paths specified")
		}

		if !updateIndexAddFlag && !updateIndexRemoveFlag {
			return fmt.Errorf("must specify either --add or --remove")
		}

		_, err := updateIndexService.UpdateIndex(args, updateIndexAddFlag, updateIndexRemoveFlag, true)
		return err
	},
}

func init() {
	updateIndexCmd.Flags().BoolVarP(&updateIndexAddFlag, "add", "a", false, "Add specified files to the index")
	updateIndexCmd.Flags().BoolVarP(
		&updateIndexRemoveFlag, "remove", "r", false, "Remove specified files from the index",
	)
	rootCmd.AddCommand(updateIndexCmd)
}
