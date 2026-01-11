package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var readTreeCmd = &cobra.Command{
	Use:   "read-tree <tree-hash>",
	Short: "Read tree objects into the index",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("tree hash required")

		}

		hash := args[0]
		return readTreeService.ReadTree(hash)

	},
}

func init() {
	rootCmd.AddCommand(readTreeCmd)
}
