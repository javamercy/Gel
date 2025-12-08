package cmd

import (
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
				cmd.PrintErrln("Error getting current working directory:", err)
				return
			}
			path = cwd
		}

		message, err := initService.Init(path)
		if err != nil {
			cmd.PrintErrln("Error initializing repository:", err)
			os.Exit(1)
		}

		cmd.Println(message)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
