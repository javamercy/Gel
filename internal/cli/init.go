package cli

import (
	"Gel/internal/gel"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new Gel repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) > 0 {
			path = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			path = cwd
		}

		initService := gel.NewInitService()
		message, err := initService.Init(path)
		if err != nil {
			return err
		}

		cmd.Println(message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
