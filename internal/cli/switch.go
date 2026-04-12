package cli

import (
	"Gel/internal/branch"

	"github.com/spf13/cobra"
)

var (
	switchCreateFlag bool
	switchForceFlag  bool
)

// switchCmd switches to an existing branch or creates and switches with --create.
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch branches or restore working tree files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		options := branch.SwitchOptions{
			Create: switchCreateFlag,
			Force:  switchForceFlag,
		}
		result, err := switchService.Switch(args[0], options)
		if err != nil {
			return err
		}

		if result.Created {
			cmd.Printf("Switched to a new branch '%s'\n", result.Branch)
		} else {
			cmd.Printf("Switched to branch '%s'\n", result.Branch)
		}
		return nil
	},
}

// init registers the switch command and its flags.
func init() {
	switchCmd.Flags().BoolVarP(
		&switchCreateFlag, "create", "c", false,
		"Create the new branch",
	)
	switchCmd.Flags().BoolVarP(
		&switchForceFlag, "force", "f", false,
		"Switch even if the index or the working tree differs from HEAD",
	)
	rootCmd.AddCommand(switchCmd)
}
