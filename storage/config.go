package storage

import (
	"Gel/core/constant"
	"Gel/core/repository"
)

type ConfigStorage struct {
	filesystemStorage IFilesystemStorage
}

func NewConfigStorage(filesystemStorage IFilesystemStorage) *ConfigStorage {
	return &ConfigStorage{
		filesystemStorage: filesystemStorage,
	}
}

func (configStorage *ConfigStorage) Read() ([]byte, error) {
	repo := repository.GetRepository()

	data, err := configStorage.filesystemStorage.ReadFile(repo.ConfigPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (configStorage *ConfigStorage) Write(data []byte) error {
	repo := repository.GetRepository()

	return configStorage.filesystemStorage.WriteFile(
		repo.ConfigPath,
		data,
		true,
		constant.GelFilePermission)
}
