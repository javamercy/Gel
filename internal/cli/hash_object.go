package cli

import (
	"github.com/spf13/cobra"
)

var (
	hashObjectWriteFlag bool
)
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object <file>...",
	Short: "Compute the hash of a file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return hashObjectService.HashObjects(cmd.OutOrStdout(), args, hashObjectWriteFlag)
	},
}

func init() {
	hashObjectCmd.Flags().BoolVarP(&hashObjectWriteFlag, "write", "w", false, "Write the object to the object database")
	rootCmd.AddCommand(hashObjectCmd)
}
