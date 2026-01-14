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
	filesystemService *vcs.FilesystemService
	objectService     *vcs.ObjectService
	indexService      *vcs.IndexService
	configService     *vcs.ConfigService

	initService        *vcs.InitService
	addService         *vcs.AddService
	hashObjectService  *vcs.HashObjectService
	catFileService     *vcs.CatFileService
	lsFilesService     *vcs.LsFilesService
	updateIndexService *vcs.UpdateIndexService
	writeTreeService   *vcs.WriteTreeService
	readTreeService    *vcs.ReadTreeService
	lsTreeService      *vcs.LsTreeService
	commitTreeService  *vcs.CommitTreeService

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

	repositoryProvider, err := repository.NewRepositoryProvider(cwd)
	if err != nil {
		return err
	}

	filesystemStorage := storage.NewFilesystemStorage()
	objectStorage := storage.NewObjectStorage(filesystemStorage, repositoryProvider)
	indexStorage := storage.NewIndexStorage(filesystemStorage, repositoryProvider)
	configStorage := storage.NewConfigStorage(filesystemStorage, repositoryProvider)

	filesystemService = vcs.NewFilesystemService(filesystemStorage)
	objectService = vcs.NewObjectService(objectStorage, filesystemService)
	indexService = vcs.NewIndexService(indexStorage)
	configService = vcs.NewConfigService(configStorage, encoding.NewBurntSushiTomlHelper())

	initService = vcs.NewInitService(filesystemService)
	hashObjectService = vcs.NewHashObjectService(objectService, filesystemService)
	catFileService = vcs.NewCatFileService(objectService)

	pathResolver := util.NewPathResolver(cwd, nil)
	updateIndexService = vcs.NewUpdateIndexService(indexService, hashObjectService, objectService)
	addService = vcs.NewAddService(updateIndexService, pathResolver)
	lsFilesService = vcs.NewLsFilesService(indexService, filesystemService, objectService)
	writeTreeService = vcs.NewWriteTreeService(indexService, objectService)
	readTreeService = vcs.NewReadTreeService(indexService, objectService)
	lsTreeService = vcs.NewLsTreeService(objectService)
	commitTreeService = vcs.NewCommitTreeService(objectService, configService)

	isServicesInitialized = true

	return nil
}
