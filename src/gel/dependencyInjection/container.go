package dependencyInjection

import (
	"Gel/src/gel/application/rules"
	"Gel/src/gel/application/services"
	"Gel/src/gel/persistence/repositories"
)

type Container struct {
	FilesystemRepository repositories.IFilesystemRepository
	ObjectRepository     repositories.IObjectRepository
	IndexRepository      repositories.IIndexRepository

	InitService        services.IInitService
	HashObjectService  services.IHashObjectService
	CatFileService     services.ICatFileService
	UpdateIndexService services.IUpdateIndexService
	LsFilesService     services.ILsFilesService
	AddService         services.IAddService
	WriteTreeService   services.IWriteTreeService

	UpdateIndexRules *rules.UpdateIndexRules
	HashObjectRules  *rules.HashObjectRules
}
