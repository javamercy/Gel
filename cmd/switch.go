package cmd

import "github.com/spf13/cobra"

var (
	switchCreateFlag bool
	switchForceFlag  bool
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch branches or restore working tree files",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		message, err := switchService.Switch(args[0], switchCreateFlag, switchForceFlag)
		if err != nil {
			return err
		}
		cmd.Println(message)
		return nil
	},
}

func init() {
	switchCmd.Flags().BoolVarP(&switchCreateFlag, "create", "c", false, "Create the new branch")
	switchCmd.Flags().BoolVarP(&switchForceFlag, "force", "f", false, "Switch even if the index or the working tree differs from HEAD")
	rootCmd.AddCommand(switchCmd)
}
