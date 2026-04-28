package cli

import (
	"Gel/internal/domain"
	"Gel/internal/inspect"

	"github.com/spf13/cobra"
)

var (
	restoreStageFlag  bool
	restoreSourceFlag string
)

// restoreCmd restores paths in the working tree or index depending on flags.
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore working tree files",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var options inspect.RestoreOptions
		switch {
		case restoreStageFlag && restoreSourceFlag != "":
			options = inspect.RestoreOptions{
				Mode:   inspect.RestoreModeCommitVsIndex,
				Source: restoreSourceFlag,
			}
		case restoreStageFlag:
			options = inspect.RestoreOptions{
				Mode: inspect.RestoreModeHEADVsIndex,
			}
		case restoreSourceFlag != "":
			options = inspect.RestoreOptions{
				Mode:   inspect.RestoreModeCommitVsWorkingTree,
				Source: restoreSourceFlag,
			}
		default:
			options = inspect.RestoreOptions{
				Mode: inspect.RestoreModeIndexVsWorkingTree,
			}
		}

		var absolutePaths []domain.AbsolutePath
		for _, arg := range args {
			absolutePath, err := domain.NewAbsolutePath(arg)
			if err != nil {
				return err
			}
			absolutePaths = append(absolutePaths, absolutePath)
		}
		return restoreService.Restore(absolutePaths, options)
	},
}

// init registers the restore command and its flags.
func init() {
	restoreCmd.Flags().BoolVarP(
		&restoreStageFlag, "staged", "s", false,
		"Restore staged files",
	)
	restoreCmd.Flags().StringVarP(
		&restoreSourceFlag, "source", "S", "",
		"Restore from specified source",
	)
	rootCmd.AddCommand(restoreCmd)
}
