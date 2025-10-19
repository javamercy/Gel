package dependencyInjection

import (
	"Gel/src/gel/application/services"
	"Gel/src/gel/core/helpers"
	"Gel/src/gel/persistence/repositories"

	"github.com/google/wire"
)

var RepositoryProviderSet = wire.NewSet(
	repositories.NewFilesystemRepository,
	wire.Bind(new(repositories.IRepository), new(*repositories.FilesystemRepository)),
)

var HelperProviderSet = wire.NewSet(
	helpers.NewZlibCompressionHelper,
	wire.Bind(new(helpers.ICompressionHelper), new(*helpers.ZlibCompressionHelper)),
)

var ServiceProviderSet = wire.NewSet(
	services.NewInitService,
	wire.Bind(new(services.IInitService), new(*services.InitService)),

	services.NewHashObjectService,
	wire.Bind(new(services.IHashObjectService), new(*services.HashObjectService)),

	services.NewCatFileService,
	wire.Bind(new(services.ICatFileService), new(*services.CatFileService)),
)
