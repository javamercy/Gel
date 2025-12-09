/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"Gel/core/repository"
	"Gel/vcs"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Core services
	filesystemService *vcs.FilesystemService
	objectService     *vcs.ObjectService
	indexService      *vcs.IndexService

	// Command services
	initService        *vcs.InitService
	addService         *vcs.AddService
	hashObjectService  *vcs.HashObjectService
	catFileService     *vcs.CatFileService
	lsFilesService     *vcs.LsFilesService
	updateIndexService *vcs.UpdateIndexService
	writeTreeService   *vcs.WriteTreeService
	readTreeService    *vcs.ReadTreeService
	lsTreeService      *vcs.LsTreeService
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gel",
	Short: "A simple version control system",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// InitializeServices sets up all the services for the commands
func InitializeServices(
	fs *vcs.FilesystemService,
	obj *vcs.ObjectService,
	idx *vcs.IndexService,
	init *vcs.InitService,
	add *vcs.AddService,
	hashObj *vcs.HashObjectService,
	catFile *vcs.CatFileService,
	lsFiles *vcs.LsFilesService,
	updateIdx *vcs.UpdateIndexService,
	writeTree *vcs.WriteTreeService,
	readTree *vcs.ReadTreeService,
	lsTree *vcs.LsTreeService,
) {
	filesystemService = fs
	objectService = obj
	indexService = idx
	initService = init
	addService = add
	hashObjectService = hashObj
	catFileService = catFile
	lsFilesService = lsFiles
	updateIndexService = updateIdx
	writeTreeService = writeTree
	readTreeService = readTree
	lsTreeService = lsTree
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Gel.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func requiresEnsureContextPreRun(cmd *cobra.Command, args []string) error {
	return repository.Initialize()
}
