package cmd

import (
	"Gel/application/services"
	"Gel/core/constants"
	"Gel/persistence/repositories"

	"github.com/spf13/cobra"
)

var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "Compute the hash of a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Please provide a file path")
			return
		}
		path := args[0]
		write, _ := cmd.Flags().GetBool("write")
		repository := repositories.NewFilesystemRepository()
		hashObjectService := services.NewHashObjectService(repository)
		hash, err := hashObjectService.HashObject(path, constants.Blob, write)
		if err != nil {
			cmd.PrintErrln("Error hashing object:", err)
			return
		}
		cmd.Println(hash)
	},
}

func init() {
	hashObjectCmd.Flags().BoolP("write", "w", false, "Write the object to the object database")
	rootCmd.AddCommand(hashObjectCmd)
}
