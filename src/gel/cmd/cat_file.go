package cmd

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/domain/objects"
	"os"

	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:     "cat-file",
	Short:   "Display the content of a Git object",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErrln("ErrorMessage: object hash required")
			_ = cmd.Help()
			os.Exit(1)
		}

		hash := args[0]

		showType, _ := cmd.Flags().GetBool("type")
		showSize, _ := cmd.Flags().GetBool("size")
		pretty, _ := cmd.Flags().GetBool("pretty")
		checkOnly, _ := cmd.Flags().GetBool("exists")

		request := dto.NewCatFileRequest(hash, showType, showSize, pretty, checkOnly)

		object, gelError := container.CatFileService.GetObject(request)
		if gelError != nil {
			cmd.PrintErrln(gelError.Message)
			os.Exit(gelError.GetExitCode())
		}

		if showType {
			cmd.Println(object.Type())
			return
		}

		if showSize {
			cmd.Println(object.Size())
			return
		}

		if object.Type() == objects.GelTreeObjectType {
			treeEntries, err := object.(*objects.Tree).DeserializeTree()
			if err != nil {
				cmd.PrintErrln(err.Error())
				os.Exit(1)
			}
			for _, entry := range treeEntries {
				objectTypeStr, ok := objects.GetObjectTypeByMode(entry.Mode)
				if !ok {
					cmd.PrintErrln("ErrorMessage: invalid object mode")
					os.Exit(1)
				}
				cmd.Printf("%s %s %s %s\n", entry.Mode, objectTypeStr, entry.Hash, entry.Name)
			}
		} else if object.Type() == objects.GelBlobObjectType {
			cmd.Println(string(object.Data()))
		}
	},
}

func init() {
	catFileCmd.Flags().BoolP("type", "t", false, "Show the object type")
	catFileCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the object content")
	catFileCmd.Flags().BoolP("size", "s", false, "Show the object size")
	catFileCmd.Flags().BoolP("exists", "e", false, "Check if the object exists")
	rootCmd.AddCommand(catFileCmd)
}
