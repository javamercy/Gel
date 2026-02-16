package cli

import "github.com/spf13/cobra"

var (
	statusShortFlag bool
)
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return statusService.Status(cmd.OutOrStdout(), statusShortFlag)
	},
}

func init() {
	statusCmd.Flags().BoolVarP(&statusShortFlag, "short", "s", false, "give the output in the short-format")
	rootCmd.AddCommand(statusCmd)
}
