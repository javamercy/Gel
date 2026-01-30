package cli

import (
	"Gel/internal/gel"
	"Gel/internal/pathspec"
	"Gel/internal/workspace"
	"Gel/storage"
	"os"

	"github.com/spf13/cobra"
)

var (
	objectService *gel.ObjectService
	indexService  *gel.IndexService
	configService *gel.ConfigService

	addService         *gel.AddService
	hashObjectService  *gel.HashObjectService
	catFileService     *gel.CatFileService
	lsFilesService     *gel.LsFilesService
	updateIndexService *gel.UpdateIndexService
	writeTreeService   *gel.WriteTreeService
	readTreeService    *gel.ReadTreeService
	lsTreeService      *gel.LsTreeService
	commitTreeService  *gel.CommitTreeService
	refService         *gel.RefService
	symbolicRefService *gel.SymbolicRefService
	updateRefService   *gel.UpdateRefService
	commitService      *gel.CommitService
	logService         *gel.LogService
	branchService      *gel.BranchService
	restoreService     *gel.RestoreService
	switchService      *gel.SwitchService

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

	filesystemStorage := storage.NewFilesystemStorage()
	objectStorage := storage.NewObjectStorage(filesystemStorage, workspaceProvider)
	indexStorage := storage.NewIndexStorage(filesystemStorage, workspaceProvider)
	configStorage := storage.NewConfigStorage(filesystemStorage, workspaceProvider)

	objectService = gel.NewObjectService(objectStorage, filesystemStorage)
	indexService = gel.NewIndexService(indexStorage)
	configService = gel.NewConfigService(configStorage)

	hashObjectService = gel.NewHashObjectService(objectService, filesystemStorage)
	catFileService = gel.NewCatFileService(objectService)

	pathResolver := pathspec.NewPathResolver(cwd, nil)
	updateIndexService = gel.NewUpdateIndexService(indexService, hashObjectService, objectService)
	addService = gel.NewAddService(indexService, updateIndexService, pathResolver)
	lsFilesService = gel.NewLsFilesService(indexService, filesystemStorage, objectService)
	writeTreeService = gel.NewWriteTreeService(indexService, objectService)
	readTreeService = gel.NewReadTreeService(indexService, objectService)
	lsTreeService = gel.NewLsTreeService(objectService)
	commitTreeService = gel.NewCommitTreeService(objectService, configService)
	refService = gel.NewRefService(workspaceProvider, filesystemStorage)
	symbolicRefService = gel.NewSymbolicRefService(refService)
	updateRefService = gel.NewUpdateRefService(refService)
	commitService = gel.NewCommitService(writeTreeService, commitTreeService, refService, filesystemStorage, objectService)
	logService = gel.NewLogService(refService, objectService)
	branchService = gel.NewBranchService(refService, workspaceProvider)
	restoreService = gel.NewRestoreService(indexService, objectService, filesystemStorage, refService)
	switchService = gel.NewSwitchService(refService, objectService, filesystemStorage, readTreeService, workspaceProvider)

	isServicesInitialized = true

	return nil
}
