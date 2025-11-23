package cmd

import (
	"Gel/src/gel/application/dto"
	"os"

	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:     "cat-file",
	Short:   "Display the content of a Git object",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErrln("Error: object hash required")
			_ = cmd.Help()
			os.Exit(1)
		}

		hash := args[0]

		showType, _ := cmd.Flags().GetBool("type")
		showSize, _ := cmd.Flags().GetBool("size")
		pretty, _ := cmd.Flags().GetBool("pretty")
		checkOnly, _ := cmd.Flags().GetBool("exists")

		request := dto.NewCatFileRequest(hash, showType, showSize, pretty, checkOnly)

		object, err := container.CatFileService.GetObject(request)
		if err != nil {
			if checkOnly {
				os.Exit(1)
			}
			cmd.PrintErrln("Error:", err)
			os.Exit(1)
		}

		if checkOnly {
			os.Exit(0)
		}

		if showType {
			cmd.Println(object.Type())
			return
		}

		if showSize {
			cmd.Println(object.Size())
			return
		}

		cmd.Print(string(object.Data()))
	},
}

func init() {
	catFileCmd.Flags().BoolP("type", "t", false, "Show the object type")
	catFileCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the object content")
	catFileCmd.Flags().BoolP("size", "s", false, "Show the object size")
	catFileCmd.Flags().BoolP("exists", "e", false, "Check if the object exists")
	rootCmd.AddCommand(catFileCmd)
}
