package dependencyInjection

import (
	"Gel/src/gel/application/services"
	"Gel/src/gel/core/helpers"
	"Gel/src/gel/persistence/repositories"
)

type Container struct {
	FilesystemRepository repositories.IFilesystemRepository
	GelRepository        repositories.IGelRepository

	CompressionHelper helpers.ICompressionHelper

	InitService       services.IInitService
	HashObjectService services.IHashObjectService
	CatFileService    services.ICatFileService
}
