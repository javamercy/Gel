package dependencyInjection

import (
	services2 "Gel/src/gel/application/services"
	"Gel/src/gel/core/helpers"
	repositories2 "Gel/src/gel/persistence/repositories"

	"github.com/google/wire"
)

var RepositoryProviderSet = wire.NewSet(
	repositories2.NewFilesystemRepository,
	wire.Bind(new(repositories2.IRepository), new(*repositories2.FilesystemRepository)),
)

var HelperProviderSet = wire.NewSet(
	helpers.NewZlibCompressionHelper,
	wire.Bind(new(helpers.ICompressionHelper), new(*helpers.ZlibCompressionHelper)),
)

var ServiceProviderSet = wire.NewSet(
	services2.NewInitService,
	wire.Bind(new(services2.IInitService), new(*services2.InitService)),

	services2.NewHashObjectService,
	wire.Bind(new(services2.IHashObjectService), new(*services2.HashObjectService)),

	services2.NewCatFileService,
	wire.Bind(new(services2.ICatFileService), new(*services2.CatFileService)),
)
