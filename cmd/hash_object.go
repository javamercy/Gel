package cmd

import (
	"github.com/spf13/cobra"
)

var (
	writeFlag bool
)
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object <file>...",
	Short: "Compute the hash of a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Error: no paths specified")
			return
		}

		write, _ := cmd.Flags().GetBool("write")

		hashMap, err := hashObjectService.HashObject(args, write)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		for _, path := range args {
			hash := hashMap[path]
			cmd.Println(hash)
		}
	},
}

func init() {
	hashObjectCmd.Flags().BoolVarP(&writeFlag, "write", "w", false, "Write the object to the object database")
	rootCmd.AddCommand(hashObjectCmd)
}
