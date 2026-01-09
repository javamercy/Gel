package cmd

import (
	"Gel/storage"
	"Gel/vcs"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new Gel repository",
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) > 0 {
			path = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			path = cwd
		}

		initService := initializeInitService()
		message, err := initService.Init(path)

		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		cmd.Println(message)
	},
}

func initializeInitService() *vcs.InitService {
	filesystemStorage := storage.NewFilesystemStorage()
	filesystemService = vcs.NewFilesystemService(filesystemStorage)
	initService = vcs.NewInitService(filesystemService)
	return initService
}

func init() {
	rootCmd.AddCommand(initCmd)
}
