package storage

import (
	"Gel/core/constant"
	"Gel/core/repository"
)

type ConfigStorage struct {
	filesystemStorage  IFilesystemStorage
	repositoryProvider repository.IRepositoryProvider
}

func NewConfigStorage(filesystemStorage IFilesystemStorage, repositoryProvider repository.IRepositoryProvider) *ConfigStorage {
	return &ConfigStorage{
		filesystemStorage:  filesystemStorage,
		repositoryProvider: repositoryProvider,
	}
}

func (configStorage *ConfigStorage) Read() ([]byte, error) {
	repo := configStorage.repositoryProvider.GetRepository()
	data, err := configStorage.filesystemStorage.ReadFile(repo.ConfigPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (configStorage *ConfigStorage) Write(data []byte) error {
	repo := configStorage.repositoryProvider.GetRepository()
	return configStorage.filesystemStorage.WriteFile(
		repo.ConfigPath,
		data,
		true,
		constant.GelFilePermission)
}
