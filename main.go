/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"Gel/cmd"
	"Gel/core/util"
	"Gel/storage"
	"Gel/vcs"
	"os"
)

func main() {
	// Storage layer
	filesystemStorage := storage.NewFilesystemStorage()
	objectStorage := storage.NewObjectStorage(filesystemStorage)
	indexStorage := storage.NewIndexStorage(filesystemStorage)

	// Core services
	filesystemService := vcs.NewFilesystemService(filesystemStorage)
	objectService := vcs.NewObjectService(objectStorage, filesystemService)
	indexService := vcs.NewIndexService(indexStorage)

	// Command services
	initService := vcs.NewInitService(objectService, filesystemService)
	hashObjectService := vcs.NewHashObjectService(objectService, filesystemService)
	catFileService := vcs.NewCatFileService(objectService)

	// Get current directory for path resolver
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	pathResolver := util.NewPathResolver(cwd)

	// Create UpdateIndexService first (needed by AddService)
	updateIndexService := vcs.NewUpdateIndexService(indexService, hashObjectService, objectService)

	// Create AddService with dependencies
	addService := vcs.NewAddService(updateIndexService, pathResolver)

	// Create remaining services
	lsFilesService := vcs.NewLsFilesService(indexService, filesystemService, objectService)
	writeTreeService := vcs.NewWriteTreeService(indexService, objectService)
	readTreeService := vcs.NewReadTreeService(indexService, objectService)
	lsTreeService := vcs.NewLsTreeService(objectService)
	commitTreeService := vcs.NewCommitTreeService(objectService)

	// Initialize commands with all services
	cmd.InitializeServices(
		filesystemService,
		objectService,
		indexService,
		initService,
		addService,
		hashObjectService,
		catFileService,
		lsFilesService,
		updateIndexService,
		writeTreeService,
		readTreeService,
		lsTreeService,
		commitTreeService,
	)

	// Execute root command
	cmd.Execute()
}
