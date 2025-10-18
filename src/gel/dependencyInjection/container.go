package dependencyInjection

import (
	services2 "Gel/src/gel/application/services"
	"Gel/src/gel/core/helpers"
	"Gel/src/gel/persistence/repositories"
)

type Container struct {
	Repository repositories.IRepository

	CompressionHelper helpers.ICompressionHelper

	InitService       services2.IInitService
	HashObjectService services2.IHashObjectService
	CatFileService    services2.ICatFileService
}
