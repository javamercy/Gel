package cmd

import "github.com/spf13/cobra"

var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "List all files tracked by Gel in the current repository",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := container.LsFilesService.LsFiles()
		if err != nil {
			cmd.PrintErrln("Error listing files:", err)
			return
		}
		for _, file := range files {
			cmd.Println(file)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsFilesCmd)
}
