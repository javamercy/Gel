package cmd

import (
	"github.com/spf13/cobra"
)

var hashObjectCmd = &cobra.Command{
	Use:     "hash-object <file>...",
	Short:   "Compute the hash of a file",
	PreRunE: requiresEnsureContextPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("Error: no paths specified")
			return
		}

		write, _ := cmd.Flags().GetBool("write")

		hashMap, _, err := hashObjectService.HashObject(args, write)
		if err != nil {
			cmd.PrintErrln("Error hashing objects:", err)
			return
		}

		for _, path := range args {
			hash := hashMap[path]
			cmd.Println(hash)
		}
	},
}

func init() {
	hashObjectCmd.Flags().BoolP("write", "w", false, "Write the object to the object database")
	hashObjectCmd.Flags().StringP("type", "t", "blob", "Specify the object type (blob, tree, commit)")
	rootCmd.AddCommand(hashObjectCmd)
}
