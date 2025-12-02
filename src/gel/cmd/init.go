package cmd

import (
	"Gel/src/gel/application/dto"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Gel repository",
	Run: func(cmd *cobra.Command, args []string) {

		var path string
		if len(args) > 0 {
			path = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				cmd.PrintErrln("Error getting current working directory:", err)
				return
			}
			path = cwd
		}

		message, gelError := container.InitService.Init(dto.NewInitRequest(path))
		if gelError != nil {
			cmd.PrintErrln(gelError.Message)
			os.Exit(gelError.GetExitCode())
		}

		cmd.Println(message)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
