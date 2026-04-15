package cli

import (
	"Gel/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	resetHardFlag  bool
	resetSoftFlag  bool
	resetMixedFlag bool
)

func init() {
	resetCmd.Flags().BoolVarP(
		&resetSoftFlag, "soft", "S", false,
		"Move HEAD only; keep index and working tree",
	)
	resetCmd.Flags().BoolVarP(
		&resetMixedFlag, "mixed", "M", false,
		"Move HEAD and reset index; keep working tree (default)",
	)
	resetCmd.Flags().BoolVarP(
		&resetHardFlag, "hard", "H", false,
		"Move HEAD, reset index, and discard working tree changes",
	)
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset [target]",
	Short: "Reset the current HEAD to a specified state",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode, err := parseResetMode()
		if err != nil {
			return err
		}

		target := ""
		if len(args) == 1 {
			target = args[0]
		}

		result, err := resetService.Reset(target, internal.ResetOptions{Mode: mode})
		if err != nil {
			return err
		}
		cmd.Printf("HEAD is now at %s\n", result.TargetHash)
		return nil
	},
}

func parseResetMode() (internal.ResetMode, error) {
	selected := 0
	mode := internal.ResetModeMixed
	if resetSoftFlag {
		selected++
		mode = internal.ResetModeSoft
	}
	if resetMixedFlag {
		selected++
		mode = internal.ResetModeMixed
	}
	if resetHardFlag {
		selected++
		mode = internal.ResetModeHard
	}
	if selected > 1 {
		return 0, fmt.Errorf("reset: only one of --soft, --mixed, --hard may be set")
	}
	return mode, nil
}
