package gel

import (
	"Gel/domain"
	"Gel/internal/storage"
	"bytes"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

const (
	ConfigSectionUser = "user"
	ConfigKeyName     = "name"
	ConfigKeyEmail    = "email"
)

type ConfigService struct {
	configStorage *storage.ConfigStorage
}

func NewConfigService(configStorage *storage.ConfigStorage) *ConfigService {
	return &ConfigService{
		configStorage: configStorage,
	}
}

func (c *ConfigService) GetUserInfo() (string, string, error) {
	config, err := c.Read()
	if err != nil {
		return "", "", fmt.Errorf("failed to decode config: %w", err)
	}
	name, ok := config.Get(ConfigSectionUser, ConfigKeyName)
	if !ok {
		return "", "", fmt.Errorf("failed to get user name: %w", err)
	}
	email, ok := config.Get(ConfigSectionUser, ConfigKeyEmail)
	if !ok {
		return "", "", fmt.Errorf("failed to get user email: %w", err)
	}
	return name, email, nil
}

func (c *ConfigService) Set(section, key, value string) error {
	config, err := c.Read()
	if err != nil {
		panic(err)
	}
	config.Set(section, key, value)
	return c.Write(config)
}

func (c *ConfigService) Get(section, key string) (string, bool) {
	config, err := c.Read()
	if err != nil {
		panic(err)
	}
	return config.Get(section, key)
}

func (c *ConfigService) List(writer io.Writer) error {
	config, err := c.Read()
	if err != nil {
		return err
	}
	for sectionName, section := range config.Sections {
		for key, value := range section {
			if _, err := fmt.Fprintf(writer, "%s.%s=%s\n", sectionName, key, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ConfigService) Write(config *domain.Config) error {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(config.Sections); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	return c.configStorage.Write(buf.Bytes())
}

func (c *ConfigService) Read() (*domain.Config, error) {
	data, err := c.configStorage.Read()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return &domain.Config{
			Sections: make(map[string]domain.Section),
		}, nil
	}

	sectionsMap := make(map[string]domain.Section)
	if _, err := toml.Decode(string(data), &sectionsMap); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}
	return domain.NewConfigFromMap(sectionsMap), nil
}
