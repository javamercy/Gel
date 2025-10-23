package cmd

import (
	"Gel/src/gel/core/constant"

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

		hash, err := container.HashObjectService.HashObject(path, constant.Blob, write)
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
