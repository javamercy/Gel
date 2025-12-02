package cmd

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/domain/objects"
	"os"

	"github.com/spf13/cobra"
)

var hashObjectCmd = &cobra.Command{
	Use:     "hash-object",
	Short:   "Compute the hash of a file",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Please provide a file path")
			return
		}
		write, _ := cmd.Flags().GetBool("write")
		objectType := cmd.Flags().Lookup("type").Value.String()

		request := dto.NewHashObjectRequest(args, objects.ObjectType(objectType), write)

		response, gelError := container.HashObjectService.HashObject(request)
		if gelError != nil {
			cmd.PrintErrln(gelError.Message)
			os.Exit(gelError.GetExitCode())
		}

		for path, hash := range response {
			cmd.Printf("%s  %s\n", hash, path)
		}
	},
}

func init() {
	hashObjectCmd.Flags().BoolP("write", "w", false, "Write the object to the object database")
	hashObjectCmd.Flags().StringP("type", "t", "blob", "Specify the object type (blob, tree, commit)")
	rootCmd.AddCommand(hashObjectCmd)
}
