package dependencyInjection

import (
	"Gel/src/gel/application/rules"
	"Gel/src/gel/application/services"
	"Gel/src/gel/persistence/repositories"

	"github.com/google/wire"
)

var PersistenceProviderSet = wire.NewSet(
	repositories.NewFilesystemRepository,
	wire.Bind(new(repositories.IFilesystemRepository), new(*repositories.FilesystemRepository)),
	repositories.NewObjectRepository,
	wire.Bind(new(repositories.IObjectRepository), new(*repositories.ObjectRepository)),
	repositories.NewIndexRepository,
	wire.Bind(new(repositories.IIndexRepository), new(*repositories.IndexRepository)),
)

var ApplicationProviderSet = wire.NewSet(
	services.NewInitService,
	wire.Bind(new(services.IInitService), new(*services.InitService)),

	rules.NewHashObjectRules,
	services.NewHashObjectService,
	wire.Bind(new(services.IHashObjectService), new(*services.HashObjectService)),

	rules.NewCatFileRules,
	services.NewCatFileService,
	wire.Bind(new(services.ICatFileService), new(*services.CatFileService)),

	services.NewUpdateIndexService,
	wire.Bind(new(services.IUpdateIndexService), new(*services.UpdateIndexService)),
	rules.NewUpdateIndexRules,

	services.NewLsFilesService,
	wire.Bind(new(services.ILsFilesService), new(*services.LsFilesService)),

	services.NewAddService,
	wire.Bind(new(services.IAddService), new(*services.AddService)),

	services.NewWriteTreeService,
	wire.Bind(new(services.IWriteTreeService), new(*services.WriteTreeService)),
)
