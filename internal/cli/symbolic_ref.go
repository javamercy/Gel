package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	symbolicRefShortFlag bool
)

// symbolicRefCmd reads or updates symbolic references such as HEAD.
var symbolicRefCmd = &cobra.Command{
	Use:   "symbolic-ref <name> [ref]",
	Short: "Read or update symbolic references",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if symbolicRefShortFlag && len(args) == 2 {
			return fmt.Errorf("symbolic-ref: --short is only valid when reading a symbolic ref")
		}

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
	symbolicRefCmd.Flags().BoolVarP(
		&symbolicRefShortFlag,
		"short", "s", false, "Shorten refs/heads/<name> to <name> when reading",
	)
	rootCmd.AddCommand(symbolicRefCmd)
}
