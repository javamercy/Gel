package cli

import "github.com/spf13/cobra"

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show various types of objects",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return showService.Show(cmd.OutOrStdout(), "")
		}
		return showService.Show(cmd.OutOrStdout(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
