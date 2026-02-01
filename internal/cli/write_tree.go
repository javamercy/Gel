package cli

import (
	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Write the current index as a tree object",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		hash, err := writeTreeService.WriteTree()
		if err != nil {
			return err
		}

		cmd.Println(hash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
