package cli

import (
	"Gel/domain"
	"Gel/internal/core"

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
		absolutePaths := make([]domain.AbsolutePath, len(args))
		for i, path := range args {
			absolutePath, err := domain.NewAbsolutePath(path)
			if err != nil {
				return err
			}
			absolutePaths[i] = absolutePath
		}
		hashes, err := hashObjectService.HashObjects(
			absolutePaths, core.HashObjectOptions{Write: hashObjectWriteFlag},
		)
		if err != nil {
			return err
		}

		for _, hash := range hashes {
			cmd.Println(hash)
		}
		return nil
	},
}

func init() {
	hashObjectCmd.Flags().BoolVarP(&hashObjectWriteFlag, "write", "w", false, "Write the object to the object database")
	rootCmd.AddCommand(hashObjectCmd)
}
