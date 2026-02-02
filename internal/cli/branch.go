package cli

import "github.com/spf13/cobra"

var (
	branchDeleteFlag bool
)
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List, create, or delete branches",
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return branchService.List(cmd.OutOrStdout())
		} else if len(args) == 1 {
			if branchDeleteFlag {
				return branchService.Delete(args[0])
			}
			return branchService.Create(args[0], "")
		}
		return branchService.Create(args[0], args[1])
	},
}

func init() {
	branchCmd.Flags().BoolVarP(&branchDeleteFlag, "delete", "d", false, "Delete branch")
	rootCmd.AddCommand(branchCmd)
}
