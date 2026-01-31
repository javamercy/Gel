package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	configListFlag bool
)

var ConfigCmd = &cobra.Command{
	Use:          "config [key] [value]",
	Short:        "Get or set repository or global options",
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if configListFlag {
			return configService.List(cmd.OutOrStdout())
		}
		if len(args) == 0 {
			return cmd.Help()
		}

		segments := strings.SplitN(args[0], ".", 2)
		if len(segments) != 2 {
			return fmt.Errorf("invalid key: %s (must be section.key)", args[0])
		}
		section, key := segments[0], segments[1]
		if len(args) == 1 {
			value, ok := configService.Get(section, key)
			if ok {
				cmd.Println(value)
			}
			return nil
		} else if len(args) == 2 {
			return configService.Set(section, key, args[1])
		}
		return nil
	},
}

func init() {
	ConfigCmd.Flags().BoolVarP(&configListFlag, "list", "l", false, "List all config values")
	rootCmd.AddCommand(ConfigCmd)
}
