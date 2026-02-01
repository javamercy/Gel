package cli

import "github.com/spf13/cobra"

var (
	updateRefDeleteFlag bool
)

var updateRefCmd = &cobra.Command{
	Use:   "update-ref",
	Short: "Update a reference",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ref := args[0]

		if updateRefDeleteFlag {
			return updateRefService.Delete(ref, args[1])
		}
		if len(args) == 2 {
			return updateRefService.Update(ref, args[1], "")
		}
		return updateRefService.Update(ref, args[1], args[2])
	},
}

func init() {
	updateRefCmd.Flags().BoolVarP(
		&updateRefDeleteFlag, "delete", "d", false, "Delete the reference instead of updating it",
	)
	rootCmd.AddCommand(updateRefCmd)
}
