package cli

import "github.com/spf13/cobra"

var (
	diffStagedFlag bool
)
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes between commits, commit and working tree, etc",
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if diffStagedFlag {
			// Head vs Index
		}
		if len(args) == 0 {
			// Index vs working dir
			return diffService.Diff()
		} else if len(args) == 1 {
			// commit vs working dir
		} else {
			// commit vs commit
		}
		return nil
	},
}

func init() {
	diffCmd.Flags().BoolVarP(&diffStagedFlag, "staged", "s", false, "Show diff between HEAD and Index")
	rootCmd.AddCommand(diffCmd)
}
