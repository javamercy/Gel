package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	dryRunFlag bool
)
var addCmd = &cobra.Command{
	Use:   "add <pathspec>...",
	Short: "Add file contents to the index",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Error: no paths specified")
			return
		}

		output, err := addService.Add(args, dryRunFlag)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cmd.Println(output)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "n", false, "Dry run the add operation without making any changes")
	rootCmd.AddCommand(addCmd)
}
