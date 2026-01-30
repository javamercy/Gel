package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	typeFlag   bool
	prettyFlag bool
	sizeFlag   bool
	existsFlag bool
)
var catFileCmd = &cobra.Command{
	Use:   "cat-file <hash>",
	Short: "Display the content of a Git object",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("hash argument required")
		}

		hash := args[0]
		return catFileService.CatFile(cmd.OutOrStdout(), hash, typeFlag, prettyFlag, sizeFlag, existsFlag)
	},
}

func init() {
	catFileCmd.Flags().BoolVarP(&typeFlag, "type", "t", false, "Show the object type")
	catFileCmd.Flags().BoolVarP(&prettyFlag, "pretty", "p", false, "Pretty-print the object content")
	catFileCmd.Flags().BoolVarP(&sizeFlag, "size", "s", false, "Show the object size")
	catFileCmd.Flags().BoolVarP(&existsFlag, "exists", "e", false, "Check if the object exists")
	rootCmd.AddCommand(catFileCmd)
}
