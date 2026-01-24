package cmd

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
			branchMap, err := branchService.List()
			if err != nil {
				return err
			}
			for name, isCurrent := range branchMap {
				if isCurrent {
					cmd.Println("* " + name)
				} else {
					cmd.Println("  " + name)
				}
			}
		} else if len(args) == 1 {
			if branchDeleteFlag {
				return branchService.Delete(args[0])
			}
			return branchService.Create(args[0])
		}
		return nil
	},
}

func init() {
	branchCmd.Flags().BoolVarP(&branchDeleteFlag, "delete", "d", false, "Delete branch")
	rootCmd.AddCommand(branchCmd)
}
