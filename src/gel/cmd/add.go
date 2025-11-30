package cmd

import (
	"Gel/src/gel/application/dto"
	"os"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add file contents to the index",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {

		all, _ := cmd.Flags().GetBool("all")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")

		var pathspecs []string
		if all {
			pathspecs = []string{"."}
		} else {
			pathspecs = args
		}

		addRequest := dto.NewAddRequest(pathspecs, dryRun, verbose)

		addResponse := container.AddService.Add(addRequest)
		if addResponse.Errors != nil {
			cmd.PrintErrln(addResponse.Errors)
			os.Exit(1)
		}

		for _, path := range addResponse.Paths {
			cmd.Println(path)
		}
	},
}

func init() {
	addCmd.Flags().BoolP("all", "A", false, "Add changes from all tracked and untracked files")
	addCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done, without making any changes")
	addCmd.Flags().BoolP("verbose", "v", false, "Show verbose output")
	rootCmd.AddCommand(addCmd)
}
