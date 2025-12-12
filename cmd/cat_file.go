package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:     "cat-file <hash>",
	Short:   "Display the content of a Git object",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErrln("Error: object hash required")
			_ = cmd.Help()
			os.Exit(1)
		}

		hash := args[0]
		objectType, _ := cmd.Flags().GetBool("type")
		pretty, _ := cmd.Flags().GetBool("pretty")
		size, _ := cmd.Flags().GetBool("size")
		exists, _ := cmd.Flags().GetBool("exists")

		output, err := catFileService.CatFile(hash, objectType, pretty, size, exists)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cmd.Println(output)
	},
}

func init() {
	catFileCmd.Flags().BoolP("type", "t", false, "Show the object type")
	catFileCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the object content")
	catFileCmd.Flags().BoolP("size", "s", false, "Show the object size")
	catFileCmd.Flags().BoolP("exists", "e", false, "Check if the object exists")
	rootCmd.AddCommand(catFileCmd)
}
