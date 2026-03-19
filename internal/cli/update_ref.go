package cli

import (
	"Gel/domain"

	"github.com/spf13/cobra"
)

var (
	updateRefDeleteFlag bool
)

var updateRefCmd = &cobra.Command{
	Use:   "update-ref",
	Short: "Update a reference",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ref := args[0]
		newHash, err := domain.NewHash(args[1])
		if err != nil {
			return err
		}
		if updateRefDeleteFlag {
			return updateRefService.Delete(ref, newHash)
		}
		if len(args) == 2 {
			return updateRefService.Update(ref, newHash, domain.Hash{})
		}
		oldHash, err := domain.NewHash(args[2])
		if err != nil {
			return err
		}
		return updateRefService.Update(ref, newHash, oldHash)
	},
}

func init() {
	updateRefCmd.Flags().BoolVarP(
		&updateRefDeleteFlag, "delete", "d", false, "Delete the reference instead of updating it",
	)
	rootCmd.AddCommand(updateRefCmd)
}
