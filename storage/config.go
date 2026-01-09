package storage

import (
	"Gel/core/constant"
	"Gel/core/repository"
)

type ConfigStorage struct {
	filesystemStorage IFilesystemStorage
	repository        *repository.Repository
}

func NewConfigStorage(filesystemStorage IFilesystemStorage, repository *repository.Repository) *ConfigStorage {
	return &ConfigStorage{
		filesystemStorage: filesystemStorage,
		repository:        repository,
	}
}

func (configStorage *ConfigStorage) Read() ([]byte, error) {

	data, err := configStorage.filesystemStorage.ReadFile(configStorage.repository.ConfigPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (configStorage *ConfigStorage) Write(data []byte) error {

	return configStorage.filesystemStorage.WriteFile(
		configStorage.repository.ConfigPath,
		data,
		true,
		constant.GelFilePermission)
}
