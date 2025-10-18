package cmd

import (
	"Gel/application/services"
	"Gel/persistence/repositories"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Gel repository",
	Run: func(cmd *cobra.Command, args []string) {
		path := "."

		if len(args) > 0 {
			path = args[0]
		}
		repository := repositories.NewFilesystemRepository()
		initService := services.NewInitService(repository)
		message, err := initService.Init(path)
		if err != nil {
			cmd.PrintErrln("Error initializing repository:", err)
			return
		}
		cmd.Println(message)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
