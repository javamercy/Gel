package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	configListFlag bool
)

// configCmd gets, sets, or lists repository config values.
var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "Get or set repository or global options",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if configListFlag {
				results, err := configService.List()
				if err != nil {
					return err
				}
				for _, result := range results {
					cmd.Println(result)
				}
				return nil
			}
			return cmd.Help()
		}

		segments := strings.SplitN(args[0], ".", 2)
		if len(segments) != 2 {
			return fmt.Errorf("invalid key: %s (must be section.key)", args[0])
		}

		section, key := segments[0], segments[1]
		if len(args) == 1 {
			value, err := configService.Get(section, key)
			if err != nil {
				return err
			}
			cmd.Println(value)
		} else if len(args) == 2 {
			return configService.Set(section, key, args[1])
		}
		return nil
	},
}

func init() {
	configCmd.Flags().BoolVarP(&configListFlag, "list", "l", false, "List all config values")
	rootCmd.AddCommand(configCmd)
}
