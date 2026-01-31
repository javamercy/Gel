package cli

import (
	"github.com/spf13/cobra"
)

var (
	lsFilesStageFlag    bool
	lsFilesCachedFlag   bool
	lsFilesDeletedFlag  bool
	lsFilesModifiedFlag bool
)
var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "List all files tracked by Gel in the current repository",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var pathspec string
		if len(args) > 0 {
			pathspec = args[0]
		}
		if !lsFilesStageFlag && !lsFilesModifiedFlag && !lsFilesDeletedFlag {
			lsFilesCachedFlag = true
		}
		return lsFilesService.LsFiles(
			cmd.OutOrStdout(), pathspec, lsFilesCachedFlag, lsFilesStageFlag,
			lsFilesModifiedFlag, lsFilesDeletedFlag,
		)
	},
}

func init() {
	lsFilesCmd.Flags().BoolVarP(
		&lsFilesCachedFlag, "cached", "c", false,
		"Show cached files in the index",
	)
	lsFilesCmd.Flags().BoolVarP(
		&lsFilesStageFlag, "stage", "s", false, "Show staged files",
	)
	lsFilesCmd.Flags().BoolVarP(
		&lsFilesModifiedFlag, "modified", "m", false, "Show modified files",
	)
	lsFilesCmd.Flags().BoolVarP(
		&lsFilesDeletedFlag, "deleted", "d", false, "Show deleted files",
	)
	rootCmd.AddCommand(lsFilesCmd)
}
