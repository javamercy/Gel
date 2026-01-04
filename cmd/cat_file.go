package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	typeFlag   bool
	prettyFlag bool
	sizeFlag   bool
	existsFlag bool
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

		output, err := catFileService.CatFile(hash, typeFlag, prettyFlag, sizeFlag, existsFlag)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cmd.Println(output)
	},
}

func init() {
	catFileCmd.Flags().BoolVarP(&typeFlag, "type", "t", false, "Show the object type")
	catFileCmd.Flags().BoolVarP(&prettyFlag, "pretty", "p", false, "Pretty-print the object content")
	catFileCmd.Flags().BoolVarP(&sizeFlag, "size", "s", false, "Show the object size")
	catFileCmd.Flags().BoolVarP(&existsFlag, "exists", "e", false, "Check if the object exists")
	rootCmd.AddCommand(catFileCmd)
}
