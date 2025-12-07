package cmd

import (
	"Gel/application/dto"
	"os"

	"github.com/spf13/cobra"
)

var lsFilesCmd = &cobra.Command{
	Use:     "ls-files",
	Short:   "List all files tracked by Gel in the current repository",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {

		stage, _ := cmd.Flags().GetBool("stage")
		cached, _ := cmd.Flags().GetBool("cached")
		deleted, _ := cmd.Flags().GetBool("deleted")
		modified, _ := cmd.Flags().GetBool("modified")

		lsFilesRequest := dto.NewLsFilesRequest(cached, stage, deleted, modified)
		files, gelError := container.LsFilesService.LsFiles(lsFilesRequest)
		if gelError != nil {
			cmd.PrintErrln(gelError)
			os.Exit(gelError.GetExitCode())
		}
		cmd.Println(files)
	},
}

func init() {
	lsFilesCmd.Flags().BoolP("stage", "s", false, "Show staged contents' mode bits, object names and stage numbers in the output")
	lsFilesCmd.Flags().BoolP("cached", "c", false, "Show staged contents' mode bits, object names and stage numbers in the output")
	lsFilesCmd.Flags().BoolP("deleted", "d", false, "Show files that have been deleted from the working directory")
	lsFilesCmd.Flags().BoolP("modified", "m", false, "Show files that have been modified in the working directory")
	rootCmd.AddCommand(lsFilesCmd)
}
