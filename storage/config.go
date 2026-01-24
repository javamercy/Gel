package storage

import (
	"Gel/core/constant"
	"Gel/core/repository"
)

type ConfigStorage struct {
	filesystemStorage  *FilesystemStorage
	repositoryProvider *repository.Provider
}

func NewConfigStorage(filesystemStorage *FilesystemStorage, repositoryProvider *repository.Provider) *ConfigStorage {
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
