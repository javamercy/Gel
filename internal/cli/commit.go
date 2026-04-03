package cli

import "github.com/spf13/cobra"

var (
	commitMessageFlag string
)

// commitCmd records the current index state as a new commit on the current branch.
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return commitService.Commit(commitMessageFlag)
	},
}

func init() {
	commitCmd.Flags().StringVarP(&commitMessageFlag, "message", "m", "", "Commit message")
	rootCmd.AddCommand(commitCmd)
}
