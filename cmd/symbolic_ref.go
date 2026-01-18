package cmd

import (
	"github.com/spf13/cobra"
)

var symbolicRefCmd = &cobra.Command{
	Use:   "symbolic-ref",
	Short: "Create, delete or list symbolic references",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if len(args) == 1 {
			ref, err := symbolicRefService.Read(name)
			if err != nil {
				return err
			}
			cmd.Println(ref)
			return nil
		}
		return symbolicRefService.Update(name, args[1])
	},
}

func init() {
	rootCmd.AddCommand(symbolicRefCmd)
}
