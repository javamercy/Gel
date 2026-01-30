package cli

import (
	"github.com/spf13/cobra"
)

var (
	stageFlag    bool
	cachedFlag   bool
	deletedFlag  bool
	modifiedFlag bool
)
var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "List all files tracked by Gel in the current repository",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !stageFlag && !modifiedFlag && !deletedFlag {
			cachedFlag = true
		}

		return lsFilesService.LsFiles(cmd.OutOrStdout(), cachedFlag, stageFlag, modifiedFlag, deletedFlag)
	},
}

func init() {
	lsFilesCmd.Flags().BoolVarP(&cachedFlag, "cached", "c", false, "Show cached files in the index")
	lsFilesCmd.Flags().BoolVarP(&stageFlag, "stage", "s", false, "Show staged files")
	lsFilesCmd.Flags().BoolVarP(&modifiedFlag, "modified", "m", false, "Show modified files")
	lsFilesCmd.Flags().BoolVarP(&deletedFlag, "deleted", "d", false, "Show deleted files")
	rootCmd.AddCommand(lsFilesCmd)
}
