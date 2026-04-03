package cli

import (
	"Gel/internal/commit"
	"Gel/internal/core"
	"Gel/internal/domain"

	"github.com/spf13/cobra"
)

var (
	logLimitFlag   int
	logOnelineFlag bool
	logSinceFlag   string
	logUntilFlag   string
)

// logCmd prints commit history starting from HEAD or a provided revision.
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := domain.HeadFileName
		if len(args) > 0 {
			name = args[0]
		}

		entries, err := logService.Log(
			name, commit.LogOptions{
				Limit:   logLimitFlag,
				Oneline: logOnelineFlag,
				Since:   logSinceFlag,
				Until:   logUntilFlag,
			},
		)
		if err != nil {
			return err
		}

		// TODO: implement paper printing

		if logOnelineFlag {
			for _, entry := range entries {
				shortHash := entry.Hash.String()[:7]
				cmd.Printf("%s%s%s %s\n", core.ColorGreen, shortHash, core.ColorReset, entry.Message)
			}
		} else {
			for _, entry := range entries {
				cmd.Printf(
					"%scommit %s%s\nDate:   %s\n\n    %s\n\n",
					core.ColorGreen,
					entry.Hash.String(),
					core.ColorReset,
					entry.Date,
					entry.Message,
				)
			}
		}
		return nil
	},
}

func init() {
	logCmd.Flags().IntVarP(
		&logLimitFlag, "limit", "n", 0,
		"Maximum number of commits to list",
	)
	logCmd.Flags().BoolVarP(
		&logOnelineFlag, "oneline", "1", false,
		"Show oneline commit summary",
	)
	logCmd.Flags().StringVarP(
		&logSinceFlag, "since", "S", "",
		"Only commits after (inclusive) this date",
	)
	logCmd.Flags().StringVarP(
		&logUntilFlag, "until", "U", "",
		"Only commits before (inclusive) this date",
	)
	rootCmd.AddCommand(logCmd)
}
