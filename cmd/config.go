package cmd

import (
	"github.com/spf13/cobra"
)

var (
	listFlag bool
)

var ConfigCmd = &cobra.Command{
	Use:     "config [key] [value]",
	Short:   "Get or set repository or global options",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {

		if listFlag {
			configMap, err := configService.List()
			if err != nil {
				cmd.PrintErrln("Error listing config:", err)
				return
			}
			for key, value := range configMap {
				cmd.Printf("%s=%s\n", key, value)
			}
			return
		}
	},
}

func init() {
	ConfigCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all config values")
	rootCmd.AddCommand(ConfigCmd)
}
