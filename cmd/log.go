package cmd

import (
	"Gel/core/constant"

	"github.com/spf13/cobra"
)

var (
	logLimitFlag   int
	logOnelineFlag bool
	logSinceFlag   string
	logUntilFlag   string
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := constant.GelHeadFileName
		if len(args) > 0 {
			name = args[0]
		}

		// TODO: handle since and until flags
		return logService.Log(cmd.OutOrStdout(), name, logLimitFlag, logOnelineFlag)
	},
}

func init() {
	logCmd.Flags().IntVarP(&logLimitFlag, "limit", "n", 0, "Maximum number of commits to list")
	logCmd.Flags().BoolVarP(&logOnelineFlag, "oneline", "1", false, "Show oneline commit summary")
	logCmd.Flags().StringVarP(&logSinceFlag, "since", "S", "", "Only commits after (inclusive) this date")
	logCmd.Flags().StringVarP(&logUntilFlag, "until", "U", "", "Only commits before (inclusive) this date")
	rootCmd.AddCommand(logCmd)
}
