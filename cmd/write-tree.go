package cmd

import (
	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Write the current index as a tree object",
	Run: func(cmd *cobra.Command, args []string) {
		hash, err := writeTreeService.WriteTree()
		if err != nil {
			cmd.PrintErrln("Error writing tree:", err)
			return
		}

		cmd.Println(hash)
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
