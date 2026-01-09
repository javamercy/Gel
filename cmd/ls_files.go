package cmd

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
	Run: func(cmd *cobra.Command, args []string) {

		if !stageFlag && !modifiedFlag && !deletedFlag {
			cachedFlag = true
		}

		output, err := lsFilesService.LsFiles(cachedFlag, stageFlag, modifiedFlag, deletedFlag)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Print(output)
	},
}

func init() {
	lsFilesCmd.Flags().BoolVarP(&cachedFlag, "cached", "c", false, "Show cached files in the index")
	lsFilesCmd.Flags().BoolVarP(&stageFlag, "stage", "s", false, "Show staged files")
	lsFilesCmd.Flags().BoolVarP(&modifiedFlag, "modified", "m", false, "Show modified files")
	lsFilesCmd.Flags().BoolVarP(&deletedFlag, "deleted", "d", false, "Show deleted files")
	rootCmd.AddCommand(lsFilesCmd)
}
