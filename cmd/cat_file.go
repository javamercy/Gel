package cmd

import (
	"Gel/domain"
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
		typeFlag, _ := cmd.Flags().GetBool("type")
		pretty, _ := cmd.Flags().GetBool("pretty")
		size, _ := cmd.Flags().GetBool("size")
		exists, _ := cmd.Flags().GetBool("exists")

		obj, err := catFileService.CatFile(hash)
		if err != nil {
			if exists {
				os.Exit(1)
			}
			cmd.PrintErrln("Error reading object:", err)
			return
		}

		if exists {
			os.Exit(0)
		}

		if typeFlag {
			cmd.Println(obj.Type())
			return
		}

		if size {
			cmd.Println(obj.Size())
			return
		}

		if pretty {
			switch o := obj.(type) {
			case *domain.Blob:
				cmd.Print(string(o.Data()))
			case *domain.Tree:
				entries, err := o.DeserializeTree()
				if err != nil {
					cmd.PrintErrln("Error deserializing tree:", err)
					return
				}
				for _, entry := range entries {
					cmd.Printf("%s %s %s\n", entry.Mode, entry.Hash, entry.Name)
				}
			default:
				cmd.PrintErrln("Unknown object type")
			}
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
