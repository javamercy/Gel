package cli

import (
	"Gel/internal/setup"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new Gel repository",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		message, err := setup.NewInitService().Init(path)
		if err != nil {
			return err
		}

		cmd.Printf("%s\n", message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
