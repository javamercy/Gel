package dependencyInjection

import (
	"Gel/src/gel/application/rules"
	"Gel/src/gel/application/services"
	"Gel/src/gel/persistence/repositories"
)

type Container struct {
	FilesystemRepository repositories.IFilesystemRepository
	GelRepository        repositories.IGelRepository

	InitService       services.IInitService
	HashObjectService services.IHashObjectService
	CatFileService    services.ICatFileService

	UpdateIndexService services.IUpdateIndexService
	UpdateIndexRules   *rules.UpdateIndexRules

	LsFilesService services.ILsFilesService
}
