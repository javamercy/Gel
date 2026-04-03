package cli

import (
	"Gel/internal/core"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	branchDeleteFlag bool
)
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List, create, or delete branches",
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if branchDeleteFlag && len(args) == 2 {
			return fmt.Errorf("branch: --delete accepts exactly one branch name")
		}
		switch len(args) {
		case 0:
			branchListItem, err := branchService.List()
			if err != nil {
				return err
			}
			for _, item := range branchListItem {
				if item.IsCurrent {
					cmd.Printf("%s* %s%s\n", core.ColorGreen, item.Name, core.ColorReset)
				} else {
					cmd.Printf("  %s\n", item.Name)
				}
			}
		case 1:
			if branchDeleteFlag {
				return branchService.Delete(args[0])
			}
			return branchService.Create(args[0], "")
		case 2:
			return branchService.Create(args[0], args[1])
		}
		return nil
	},
}

func init() {
	branchCmd.Flags().BoolVarP(
		&branchDeleteFlag, "delete", "d", false,
		"Delete branch",
	)
	rootCmd.AddCommand(branchCmd)
}
