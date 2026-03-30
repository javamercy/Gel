package cli

import (
	"Gel/domain"
	"Gel/internal/staging"

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
		absolutePaths := make([]domain.AbsolutePath, len(args))
		for i, path := range args {
			absolutePath, err := domain.NewAbsolutePath(path)
			if err != nil {
				return err
			}
			absolutePaths[i] = absolutePath
		}

		_, err := updateIndexService.UpdateIndex(
			absolutePaths,
			staging.UpdateIndexOptions{
				Add:    updateIndexAddFlag,
				Remove: updateIndexRemoveFlag,
				Write:  true,
			},
		)
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
