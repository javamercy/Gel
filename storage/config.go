package storage

import (
	"Gel/core/encoding"
	"Gel/core/repository"
	"Gel/domain"
)

type ConfigStorage struct {
	filesystemStorage IFilesystemStorage
	toml              encoding.IToml
}

func NewConfigStorage(filesystemStorage IFilesystemStorage) *ConfigStorage {
	return &ConfigStorage{
		filesystemStorage: filesystemStorage,
		toml:              encoding.NewBurntSushiToml(),
	}
}

func (configStorage *ConfigStorage) Read() (*domain.Config, error) {
	repo := repository.GetRepository()
	var config *domain.Config

	data, err := configStorage.filesystemStorage.ReadFile(repo.ConfigPath)
	if err != nil {
		return nil, err
	}

	if err := configStorage.toml.Decode(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func (configStorage *ConfigStorage) Write(config *domain.Config) error {
	repo := repository.GetRepository()
	encodedData, err := configStorage.toml.Encode(config)
	if err != nil {
		return err
	}

	return configStorage.filesystemStorage.WriteFile(
		repo.ConfigPath,
		encodedData,
		true,
		0644)
}
