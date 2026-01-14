package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
	"Gel/storage"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrUnknownConfigKey = errors.New("unknown config key")
	ErrInvalidKeyFormat = errors.New("invalid key format")
)

const (
	ConfigSectionUserStr = "user"
	ConfigKeyNameStr     = "name"
	ConfigKeyEmailStr    = "email"
)

type ConfigService struct {
	configStorage *storage.ConfigStorage
	tomlHelper    encoding.ITomlHelper
}

func NewConfigService(configStorage *storage.ConfigStorage, tomlHelper encoding.ITomlHelper) *ConfigService {
	return &ConfigService{
		configStorage: configStorage,
		tomlHelper:    tomlHelper,
	}
}

func (configService *ConfigService) DecodeConfig() (*domain.Config, error) {
	config := domain.Config{}
	data, err := configService.configStorage.Read()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return &config, nil
	}

	if err := configService.tomlHelper.Decode(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (configService *ConfigService) GetUserIdentity() (domain.UserIdentity, error) {

	var user domain.UserIdentity
	config, err := configService.DecodeConfig()

	if err != nil {
		return user, err
	}
	return domain.NewUserIdentity(config.User.Name, config.User.Email)
}

func (configService *ConfigService) Get(key string) (string, error) {
	config, err := configService.DecodeConfig()
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(key, ".", 2)
	if len(parts) != 2 {
		return "", ErrInvalidKeyFormat
	}

	section, keyName := parts[0], parts[1]

	switch section {
	case ConfigSectionUserStr:
		switch keyName {
		case ConfigKeyNameStr:
			return config.User.Name, nil
		case ConfigKeyEmailStr:
			return config.User.Email, nil
		}
	}

	return "", ErrUnknownConfigKey
}

func (configService *ConfigService) Set(key, value string) error {
	config, err := configService.DecodeConfig()
	if err != nil {
		return err
	}

	parts := strings.SplitN(key, ".", 2)
	if len(parts) != 2 {
		return ErrInvalidKeyFormat
	}

	section, keyName := parts[0], parts[1]

	switch section {
	case ConfigSectionUserStr:
		switch keyName {
		case ConfigKeyNameStr:
			config.User.Name = value
		case ConfigKeyEmailStr:
			config.User.Email = value
		}
	default:
		return ErrUnknownConfigKey
	}

	data, err := configService.tomlHelper.Encode(config)
	if err != nil {
		return err
	}
	return configService.configStorage.Write(data)
}

func (configService *ConfigService) List(writer io.Writer) error {
	config, err := configService.DecodeConfig()
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(writer,
		"%s.%s=%s\n"+
			"%s.%s=%s\n",
		ConfigSectionUserStr, ConfigKeyNameStr, config.User.Name,
		ConfigSectionUserStr, ConfigKeyEmailStr, config.User.Email); err != nil {
		return err
	}
	return nil
}
