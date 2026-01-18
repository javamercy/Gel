package cmd

import "github.com/spf13/cobra"

var updateRefCmd = &cobra.Command{
	Use:   "update-ref",
	Short: "Update a reference",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {

		ref := args[0]
		newHash := args[1]
		if len(args) == 2 {
			return updateRefService.Update(ref, newHash)
		}

		oldHash := args[2]
		return updateRefService.SafeUpdate(ref, newHash, oldHash)
	},
}

func init() {
	rootCmd.AddCommand(updateRefCmd)
}
