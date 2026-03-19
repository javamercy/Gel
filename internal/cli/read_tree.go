package cli

import (
	"Gel/domain"

	"github.com/spf13/cobra"
)

var readTreeCmd = &cobra.Command{
	Use:   "read-tree <tree-hash>",
	Short: "Read tree objects into the index",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash, err := domain.NewHash(args[0])
		if err != nil {
			return err
		}
		return readTreeService.ReadTree(hash)

	},
}

func init() {
	rootCmd.AddCommand(readTreeCmd)
}
