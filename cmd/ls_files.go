package cmd

import (
	"github.com/spf13/cobra"
)

var lsFilesCmd = &cobra.Command{
	Use:     "ls-files",
	Short:   "List all files tracked by Gel in the current repository",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		cached, _ := cmd.Flags().GetBool("cached")
		stage, _ := cmd.Flags().GetBool("stage")
		modified, _ := cmd.Flags().GetBool("modified")
		deleted, _ := cmd.Flags().GetBool("deleted")

		if !stage && !modified && !deleted {
			cached = true
		}

		output, err := lsFilesService.LsFiles(cached, stage, modified, deleted)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Print(output)
	},
}

func init() {
	lsFilesCmd.Flags().BoolP("stage", "s", false, "Show staged contents' mode bits, object names and stage numbers in the output")
	lsFilesCmd.Flags().BoolP("cached", "c", false, "Show staged contents' mode bits, object names and stage numbers in the output")
	lsFilesCmd.Flags().BoolP("deleted", "d", false, "Show files that have been deleted from the working directory")
	lsFilesCmd.Flags().BoolP("modified", "m", false, "Show files that have been modified in the working directory")
	rootCmd.AddCommand(lsFilesCmd)
}
