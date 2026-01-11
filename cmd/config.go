package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listFlag bool
)

var ConfigCmd = &cobra.Command{
	Use:          "config [key] [value]",
	Short:        "Get or set repository or global options",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if listFlag {
			return configService.List(cmd.OutOrStdout())
		}

		if len(args) == 1 {
			value, err := configService.Get(args[0])
			if err != nil {
				return err
			}

			cmd.Println(value)
			return nil
		}

		if len(args) == 2 {
			return configService.Set(args[0], args[1])
		}

		return fmt.Errorf("invalid usage: use 'gel config --list' or 'gel config <key>' or 'gel config <key> <value>'")
	},
}

func init() {
	ConfigCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all config values")
	rootCmd.AddCommand(ConfigCmd)
}
