package cli

import (
	"Gel/internal/branch"
	"Gel/internal/commit"
	"Gel/internal/core"
	"Gel/internal/diff"
	"Gel/internal/inspect"
	"Gel/internal/staging"
	"Gel/internal/storage"
	"Gel/internal/tree"
	"Gel/internal/workspace"
	"os"

	"github.com/spf13/cobra"
)

var (
	objectService     *core.ObjectService
	indexService      *core.IndexService
	configService     *core.ConfigService
	refService        *core.RefService
	hashObjectService *core.HashObjectService
	treeResolver      *core.TreeResolver
	pathResolver      *core.PathResolver
	changeDetector    *core.ChangeDetector
)

var (
	addService         *staging.AddService
	catFileService     *inspect.CatFileService
	lsFilesService     *staging.LsFilesService
	updateIndexService *staging.UpdateIndexService
	writeTreeService   *tree.WriteTreeService
	readTreeService    *tree.ReadTreeService
	lsTreeService      *tree.LsTreeService
	commitTreeService  *commit.CommitTreeService
	symbolicRefService *core.SymbolicRefService
	updateRefService   *core.UpdateRefService
	commitService      *commit.CommitService
	logService         *commit.LogService
	branchService      *branch.BranchService
	restoreService     *inspect.RestoreService
	switchService      *branch.SwitchService
	statusService      *inspect.StatusService
	diffService        *diff.DiffService
	showService        *inspect.ShowService

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

	objectService = core.NewObjectService(objectStorage)
	indexService = core.NewIndexService(indexStorage)
	configService = core.NewConfigService(configStorage)
	refService = core.NewRefService(workspaceProvider)
	hashObjectService = core.NewHashObjectService(objectService)
	pathResolver = core.NewPathResolver(cwd, nil)
	changeDetector = core.NewChangeDetector(hashObjectService)
	treeResolver = core.NewTreeResolver(
		objectService, indexService, refService, pathResolver, hashObjectService, changeDetector,
	)
	symbolicRefService = core.NewSymbolicRefService(refService)
	updateRefService = core.NewUpdateRefService(refService)

	catFileService = inspect.NewCatFileService(objectService)
	updateIndexService = staging.NewUpdateIndexService(indexService, objectService, hashObjectService, changeDetector)
	addService = staging.NewAddService(indexService, updateIndexService, pathResolver)
	lsFilesService = staging.NewLsFilesService(indexService, objectService, changeDetector)
	writeTreeService = tree.NewWriteTreeService(indexService, objectService)
	readTreeService = tree.NewReadTreeService(indexService, objectService)
	lsTreeService = tree.NewLsTreeService(objectService)
	commitTreeService = commit.NewCommitTreeService(objectService, configService)
	commitService = commit.NewCommitService(writeTreeService, commitTreeService, refService, objectService)
	logService = commit.NewLogService(refService, objectService)
	switchService = branch.NewSwitchService(
		indexService, refService, branchService, objectService, readTreeService, treeResolver, restoreService,
	)
	branchService = branch.NewBranchService(refService, objectService, workspaceProvider)
	restoreService = inspect.NewRestoreService(indexService, objectService, refService, treeResolver, changeDetector)
	statusService = inspect.NewStatusService(indexService, objectService, treeResolver, refService, symbolicRefService)
	diffService = diff.NewDiffService(objectService, refService, treeResolver, diff.NewMyersDiffAlgorithm())
	showService = inspect.NewShowService(objectService, refService, diffService)

	isServicesInitialized = true
	return nil
}
