package cmd

import (
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

		message, err := container.InitService.Init(path)
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
