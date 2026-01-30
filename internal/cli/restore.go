package cli

import "github.com/spf13/cobra"

var (
	restoreStageFlag  bool
	restoreSourceFlag string
)
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore working tree files",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return restoreService.Restore(args, restoreSourceFlag, restoreStageFlag)
	},
}

func init() {
	restoreCmd.Flags().BoolVarP(&restoreStageFlag, "staged", "s", false, "Restore staged files")
	restoreCmd.Flags().StringVarP(&restoreSourceFlag, "source", "S", "", "Restore from specified source")
	rootCmd.AddCommand(restoreCmd)
}
