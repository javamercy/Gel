package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	writeFlag bool
)
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object <file>...",
	Short: "Compute the hash of a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no paths specified")
		}
		return hashObjectService.HashObjects(cmd.OutOrStdout(), args, writeFlag)
	},
}

func init() {
	hashObjectCmd.Flags().BoolVarP(&writeFlag, "write", "w", false, "Write the object to the object database")
	rootCmd.AddCommand(hashObjectCmd)
}
