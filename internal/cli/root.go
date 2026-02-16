package cli

import (
	"Gel/internal"
	branch2 "Gel/internal/branch"
	commit2 "Gel/internal/commit"
	core2 "Gel/internal/core"
	diff2 "Gel/internal/diff"
	inspect2 "Gel/internal/inspect"
	staging2 "Gel/internal/staging"
	"Gel/internal/storage"
	tree2 "Gel/internal/tree"
	"Gel/internal/workspace"
	"os"

	"github.com/spf13/cobra"
)

var (
	objectService     *core2.ObjectService
	indexService      *core2.IndexService
	configService     *core2.ConfigService
	refService        *core2.RefService
	hashObjectService *core2.HashObjectService
	treeResolver      *core2.TreeResolver
	pathResolver      *core2.PathResolver
)

var (
	addService            *staging2.AddService
	catFileService        *inspect2.CatFileService
	lsFilesService        *staging2.LsFilesService
	updateIndexService    *staging2.UpdateIndexService
	writeTreeService      *tree2.WriteTreeService
	readTreeService       *tree2.ReadTreeService
	lsTreeService         *tree2.LsTreeService
	commitTreeService     *commit2.CommitTreeService
	symbolicRefService    *internal.SymbolicRefService
	updateRefService      *internal.UpdateRefService
	commitService         *commit2.CommitService
	logService            *commit2.LogService
	branchService         *branch2.BranchService
	restoreService        *inspect2.RestoreService
	switchService         *branch2.SwitchService
	statusService         *inspect2.StatusService
	diffService           *diff2.DiffService
	isServicesInitialized bool
)

var commandsWithoutRepository = map[string]bool{
	"init": true,
	"help": true,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gel",
	Short: "An Agentic Version Control System",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if commandsWithoutRepository[cmd.Name()] {
			return nil
		}
		return initializeServices()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initializeServices sets up all services lazily when a command needs them
func initializeServices() error {
	if isServicesInitialized {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	workspaceProvider, err := workspace.NewProvider(cwd)
	if err != nil {
		return err
	}

	objectStorage := storage.NewObjectStorage(workspaceProvider)
	indexStorage := storage.NewIndexStorage(workspaceProvider)
	configStorage := storage.NewConfigStorage(workspaceProvider)

	objectService = core2.NewObjectService(objectStorage)
	indexService = core2.NewIndexService(indexStorage)
	configService = core2.NewConfigService(configStorage)
	refService = core2.NewRefService(workspaceProvider)
	hashObjectService = core2.NewHashObjectService(objectService)
	pathResolver = core2.NewPathResolver(cwd, nil)
	treeResolver = core2.NewTreeResolver(objectService, indexService, refService, pathResolver, hashObjectService)

	catFileService = inspect2.NewCatFileService(objectService)
	updateIndexService = staging2.NewUpdateIndexService(indexService, hashObjectService, objectService)
	addService = staging2.NewAddService(indexService, updateIndexService, pathResolver)
	lsFilesService = staging2.NewLsFilesService(indexService, objectService)
	writeTreeService = tree2.NewWriteTreeService(indexService, objectService)
	readTreeService = tree2.NewReadTreeService(indexService, objectService)
	lsTreeService = tree2.NewLsTreeService(objectService)
	commitTreeService = commit2.NewCommitTreeService(objectService, configService)

	symbolicRefService = internal.NewSymbolicRefService(refService)
	updateRefService = internal.NewUpdateRefService(refService)
	commitService = commit2.NewCommitService(writeTreeService, commitTreeService, refService, objectService)
	logService = commit2.NewLogService(refService, objectService)
	branchService = branch2.NewBranchService(refService, objectService, workspaceProvider)
	restoreService = inspect2.NewRestoreService(indexService, objectService, hashObjectService, refService)
	switchService = branch2.NewSwitchService(refService, objectService, readTreeService, workspaceProvider)
	statusService = inspect2.NewStatusService(indexService, objectService, treeResolver, refService, symbolicRefService)
	diffService = diff2.NewDiffService(objectService, refService, treeResolver, diff2.NewMyersDiffAlgorithm())
	isServicesInitialized = true
	return nil
}
