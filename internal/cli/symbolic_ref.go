package cli

import (
	"github.com/spf13/cobra"
)

var (
	symbolicRefShortFlag bool
)

var symbolicRefCmd = &cobra.Command{
	Use:   "symbolic-ref",
	Short: "Create, delete or list symbolic references",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if len(args) == 1 {
			ref, err := symbolicRefService.Read(name, symbolicRefShortFlag)
			if err != nil {
				return err
			}
			cmd.Println(ref)
			return nil
		}
		return symbolicRefService.Write(name, args[1])
	},
}

func init() {
	symbolicRefCmd.Flags().BoolVarP(&symbolicRefShortFlag, "short", "s", false, "Print the short name of the reference")
	rootCmd.AddCommand(symbolicRefCmd)
}
