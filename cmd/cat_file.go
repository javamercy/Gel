package cmd

import (
	"Gel/application/services"
	"Gel/persistence/repositories"
	"os"

	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Display the content of a Git object",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_ = cmd.Help()
			os.Exit(1)
		}
		hash := args[0]
		repository := repositories.NewFilesystemRepository()
		catFileService := services.NewCatFileService(repository)

		exists, _ := cmd.Flags().GetBool("exists")

		if exists {
			exists := catFileService.ObjectExists(hash)
			if exists {
				cmd.Println("Object exists")
				os.Exit(0)
			}
		}

		object, err := catFileService.GetObject(hash)
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		cmd.Println("Object:", string(object.Data()))

	},
}

// -t (type): Displays the object type (e.g., "blob", "tree", "commit"). If the object doesn't exist, it errors.
// -s (size): Displays the size of the object's content in bytes (not including the header).
// -p (pretty-print): Pretty-prints the object's content. For blobs, this is just the raw content. For trees/commits, it's formatted (e.g., tree entries listed, commit details shown). This is the most common flag.
// -e (exists): Exits with status 0 if the object exists and is valid, non-zero otherwise. Often used in scripts for validation.
// Object Argument: The hash (full or partial) of the object to inspect. Git resolves partial hashes if unique.

func init() {
	catFileCmd.Flags().StringP("type", "t", "", "Specify the type of the object (e.g., blob, tree, commit)")
	catFileCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the object content")
	catFileCmd.Flags().BoolP("size", "s", false, "Display the size of the object")
	catFileCmd.Flags().BoolP("exists", "e", false, "Check if the object exists")
	rootCmd.AddCommand(catFileCmd)
}
