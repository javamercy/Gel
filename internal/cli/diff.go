package cli

import (
	"Gel/domain"
	"Gel/internal/diff"
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
			baseCommitHash, err := domain.NewHash(arg)
			if err != nil {
				return err
			}
			return diffService.Diff(
				cmd.OutOrStdout(), diff.DiffOptions{Mode: diff.ModeCommitVsWorkingTree, BaseCommitHash: baseCommitHash},
			)
		}

		baseCommitHash, err := domain.NewHash(args[0])
		if err != nil {
			return err
		}
		targetCommitHash, err := domain.NewHash(args[1])
		if err != nil {
			return err
		}
		return diffService.Diff(
			cmd.OutOrStdout(),
			diff.DiffOptions{
				Mode: diff.ModeCommitVsCommit, BaseCommitHash: baseCommitHash, TargetCommitHash: targetCommitHash,
			},
		)
	},
}

func init() {
	diffCmd.Flags().BoolVarP(&diffStagedFlag, "staged", "s", false, "Show diff between HEAD and Index")
	rootCmd.AddCommand(diffCmd)
}
