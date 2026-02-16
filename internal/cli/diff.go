package cli

import (
	"Gel/internal/gel/diff"
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
			if diffStagedFlag {
				return diffService.Diff(cmd.OutOrStdout(), diff.DiffOptions{Mode: diff.ModeIndexVsHEAD})
			}
			return diffService.Diff(cmd.OutOrStdout(), diff.DiffOptions{Mode: diff.ModeWorkingTreeVsIndex})
		} else if len(args) == 1 {
			arg := args[0]
			if arg == workspace.HeadFileName {
				return diffService.Diff(cmd.OutOrStdout(), diff.DiffOptions{Mode: diff.ModeWorkingTreeVsHEAD})
			}
			return diffService.Diff(
				cmd.OutOrStdout(), diff.DiffOptions{Mode: diff.ModeCommitVsWorkingTree, BaseCommitHash: arg},
			)
		}
		return diffService.Diff(
			cmd.OutOrStdout(),
			diff.DiffOptions{Mode: diff.ModeCommitVsCommit, BaseCommitHash: args[0], TargetCommitHash: args[1]},
		)
	},
}

func init() {
	diffCmd.Flags().BoolVarP(&diffStagedFlag, "staged", "s", false, "Show diff between HEAD and Index")
	rootCmd.AddCommand(diffCmd)
}
