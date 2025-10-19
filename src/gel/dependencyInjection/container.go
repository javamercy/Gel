package dependencyInjection

import (
	"Gel/src/gel/application/services"
	"Gel/src/gel/core/helpers"
	"Gel/src/gel/persistence/repositories"
)

type Container struct {
	Repository repositories.IRepository

	CompressionHelper helpers.ICompressionHelper

	InitService       services.IInitService
	HashObjectService services.IHashObjectService
	CatFileService    services.ICatFileService
}
