package cmd

import "github.com/spf13/cobra"

var (
	commitMessageFlag string
)
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
