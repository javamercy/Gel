package cli

import (
	"Gel/internal/workspace"

	"github.com/spf13/cobra"
)

var (
	diffStagedFlag bool
)
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes between commits, commit and working tree, etc",
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return diffService.Diff(false, diffStagedFlag, "", "")
		} else if len(args) == 1 {
			arg := args[0]
			if arg == workspace.HeadFileName {
				return diffService.Diff(true, false, "", "")
			}
			return diffService.Diff(false, false, arg, "")
		}
		return diffService.Diff(false, false, args[0], args[1])
	},
}

func init() {
	diffCmd.Flags().BoolVarP(&diffStagedFlag, "staged", "s", false, "Show diff between HEAD and Index")
	rootCmd.AddCommand(diffCmd)
}
