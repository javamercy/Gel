package cli

import (
	"Gel/internal/core"

	"github.com/spf13/cobra"
)

// statusCmd shows staged, unstaged, and untracked working tree changes.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := statusService.Status()
		if err != nil {
			return err
		}

		cmd.Printf("On branch %s%s%s\n", core.ColorGreen, result.CurrentBranch, core.ColorReset)

		if result.HeadTreeSize == 0 {
			cmd.Println("No commits yet")
		}
		if len(result.Staged) > 0 {
			cmd.Printf("\n%sChanges to be committed:%s\n", core.ColorGreen, core.ColorReset)
			for _, staged := range result.Staged {
				cmd.Printf(
					"\t%s%s:  %s%s\n",
					core.ColorGreen, staged.Status, staged.Path, core.ColorReset,
				)
			}
		}
		if len(result.Unstaged) > 0 {
			cmd.Printf("\n%sChanges not staged for commit:%s\n", core.ColorRed, core.ColorReset)
			for _, unstaged := range result.Unstaged {
				cmd.Printf(
					"\t%s:  %s%s\n",
					unstaged.Status, unstaged.Path, core.ColorReset,
				)
			}
		}
		if len(result.Untracked) > 0 {
			cmd.Printf("\n%sUntracked files:%s\n", core.ColorRed, core.ColorReset)
			for _, untracked := range result.Untracked {
				cmd.Printf("\t%s%s\n", untracked, core.ColorReset)
			}
		}
		return nil
	},
}

// init registers the status command.
func init() {
	rootCmd.AddCommand(statusCmd)
}
