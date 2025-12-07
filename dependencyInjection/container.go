package dependencyInjection

import (
	rules2 "Gel/application/rules"
	services2 "Gel/application/services"
	repositories2 "Gel/persistence/repositories"
)

type Container struct {
	FilesystemRepository repositories2.IFilesystemRepository
	ObjectRepository     repositories2.IObjectRepository
	IndexRepository      repositories2.IIndexRepository

	InitService        services2.IInitService
	HashObjectService  services2.IHashObjectService
	CatFileService     services2.ICatFileService
	UpdateIndexService services2.IUpdateIndexService
	LsFilesService     services2.ILsFilesService
	AddService         services2.IAddService
	WriteTreeService   services2.IWriteTreeService
	ReadTreeService    services2.IReadTreeService

	UpdateIndexRules *rules2.UpdateIndexRules
	HashObjectRules  *rules2.HashObjectRules
}
