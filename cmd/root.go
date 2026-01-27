package cmd

import (
	"Gel/core/encoding"
	"Gel/core/repository"
	"Gel/core/util"
	"Gel/storage"
	"Gel/vcs"
	"os"

	"github.com/spf13/cobra"
)

var (
	objectService *vcs.ObjectService
	indexService  *vcs.IndexService
	configService *vcs.ConfigService

	addService         *vcs.AddService
	hashObjectService  *vcs.HashObjectService
	catFileService     *vcs.CatFileService
	lsFilesService     *vcs.LsFilesService
	updateIndexService *vcs.UpdateIndexService
	writeTreeService   *vcs.WriteTreeService
	readTreeService    *vcs.ReadTreeService
	lsTreeService      *vcs.LsTreeService
	commitTreeService  *vcs.CommitTreeService
	refService         *vcs.RefService
	symbolicRefService *vcs.SymbolicRefService
	updateRefService   *vcs.UpdateRefService
	commitService      *vcs.CommitService
	logService         *vcs.LogService
	branchService      *vcs.BranchService
	restoreService     *vcs.RestoreService

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

	repositoryProvider, err := repository.NewProvider(cwd)
	if err != nil {
		return err
	}

	filesystemStorage := storage.NewFilesystemStorage()
	objectStorage := storage.NewObjectStorage(filesystemStorage, repositoryProvider)
	indexStorage := storage.NewIndexStorage(filesystemStorage, repositoryProvider)
	configStorage := storage.NewConfigStorage(filesystemStorage, repositoryProvider)

	objectService = vcs.NewObjectService(objectStorage, filesystemStorage)
	indexService = vcs.NewIndexService(indexStorage)
	configService = vcs.NewConfigService(configStorage, encoding.NewBurntSushiTomlHelper())

	hashObjectService = vcs.NewHashObjectService(objectService, filesystemStorage)
	catFileService = vcs.NewCatFileService(objectService)

	pathResolver := util.NewPathResolver(cwd, nil)
	updateIndexService = vcs.NewUpdateIndexService(indexService, hashObjectService, objectService)
	addService = vcs.NewAddService(updateIndexService, pathResolver)
	lsFilesService = vcs.NewLsFilesService(indexService, filesystemStorage, objectService)
	writeTreeService = vcs.NewWriteTreeService(indexService, objectService)
	readTreeService = vcs.NewReadTreeService(indexService, objectService)
	lsTreeService = vcs.NewLsTreeService(objectService)
	commitTreeService = vcs.NewCommitTreeService(objectService, configService)
	refService = vcs.NewRefService(repositoryProvider, filesystemStorage)
	symbolicRefService = vcs.NewSymbolicRefService(refService)
	updateRefService = vcs.NewUpdateRefService(refService)
	commitService = vcs.NewCommitService(
		writeTreeService,
		commitTreeService,
		refService,
		filesystemStorage,
		objectService)
	logService = vcs.NewLogService(refService, objectService)
	branchService = vcs.NewBranchService(refService, repositoryProvider)
	restoreService = vcs.NewRestoreService(indexService, objectService, filesystemStorage, refService)

	isServicesInitialized = true

	return nil
}
