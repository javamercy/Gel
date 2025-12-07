package dependencyInjection

import (
	rules2 "Gel/application/rules"
	services2 "Gel/application/services"
	repositories2 "Gel/persistence/repositories"

	"github.com/google/wire"
)

var PersistenceProviderSet = wire.NewSet(
	repositories2.NewFilesystemRepository,
	wire.Bind(new(repositories2.IFilesystemRepository), new(*repositories2.FilesystemRepository)),
	repositories2.NewObjectRepository,
	wire.Bind(new(repositories2.IObjectRepository), new(*repositories2.ObjectRepository)),
	repositories2.NewIndexRepository,
	wire.Bind(new(repositories2.IIndexRepository), new(*repositories2.IndexRepository)),
)

var ApplicationProviderSet = wire.NewSet(
	services2.NewInitService,
	wire.Bind(new(services2.IInitService), new(*services2.InitService)),

	rules2.NewHashObjectRules,
	services2.NewHashObjectService,
	wire.Bind(new(services2.IHashObjectService), new(*services2.HashObjectService)),

	rules2.NewCatFileRules,
	services2.NewCatFileService,
	wire.Bind(new(services2.ICatFileService), new(*services2.CatFileService)),

	services2.NewUpdateIndexService,
	wire.Bind(new(services2.IUpdateIndexService), new(*services2.UpdateIndexService)),
	rules2.NewUpdateIndexRules,

	services2.NewLsFilesService,
	wire.Bind(new(services2.ILsFilesService), new(*services2.LsFilesService)),

	services2.NewAddService,
	wire.Bind(new(services2.IAddService), new(*services2.AddService)),

	services2.NewWriteTreeService,
	wire.Bind(new(services2.IWriteTreeService), new(*services2.WriteTreeService)),

	services2.NewReadTreeService,
	wire.Bind(new(services2.IReadTreeService), new(*services2.ReadTreeService)),
)
