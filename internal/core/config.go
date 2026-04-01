package core

import (
	"Gel/internal/domain"
	"Gel/internal/storage"
	"Gel/internal/validate"
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
		return "", "", fmt.Errorf("config: failed to read: %w", err)
	}

	name, ok := config.Get(ConfigSectionUser, ConfigKeyName)
	if !ok {
		return "", "", fmt.Errorf("config: user.name is not set")
	}

	email, ok := config.Get(ConfigSectionUser, ConfigKeyEmail)
	if !ok {
		return "", "", fmt.Errorf("config: user.email is not set")
	}
	return name, email, nil
}

func (c *ConfigService) Set(section, key, value string) error {
	if err := validate.StringMustNotBeEmpty(section); err != nil {
		return fmt.Errorf("config: %w", err)
	}
	if err := validate.StringMustNotBeEmpty(key); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	config, err := c.Read()
	if err != nil {
		return fmt.Errorf("config: failed to read: %w", err)
	}
	config.Set(section, key, value)
	return c.Write(config)
}

func (c *ConfigService) GetAndOutput(writer io.Writer, section, key string) error {
	value, err := c.Get(section, key)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintln(writer, value); err != nil {
		return fmt.Errorf("config: failed to write config: %w", err)
	}
	return nil
}

func (c *ConfigService) Get(section, key string) (string, error) {
	if err := validate.StringMustNotBeEmpty(section); err != nil {
		return "", fmt.Errorf("config: %w", err)
	}
	if err := validate.StringMustNotBeEmpty(key); err != nil {
		return "", fmt.Errorf("config: %w", err)
	}

	config, err := c.Read()
	if err != nil {
		return "", fmt.Errorf("config: failed to read: %w", err)
	}
	value, ok := config.Get(section, key)
	if !ok {
		return "", fmt.Errorf("config: key '%s.%s' not found", section, key)
	}
	return value, nil
}

func (c *ConfigService) List(writer io.Writer) error {
	config, err := c.Read()
	if err != nil {
		return err
	}
	for sectionName, section := range config.Sections {
		for key, value := range section {
			if _, err := fmt.Fprintf(writer, "%s.%s=%s\n", sectionName, key, value); err != nil {
				return fmt.Errorf("config: failed to write config: %w", err)
			}
		}
	}
	return nil
}

func (c *ConfigService) Write(config *domain.Config) error {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(config.Sections); err != nil {
		return fmt.Errorf("config: failed to encode config: %w", err)
	}
	return c.configStorage.Write(buf.Bytes())
}

func (c *ConfigService) Read() (*domain.Config, error) {
	data, err := c.configStorage.Read()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if len(data) == 0 {
		return &domain.Config{
			Sections: make(map[string]domain.Section),
		}, nil
	}

	sectionsMap := make(map[string]domain.Section)
	if _, err := toml.Decode(string(data), &sectionsMap); err != nil {
		return nil, fmt.Errorf("config: failed to decode config: %w", err)
	}
	return domain.NewConfigFromMap(sectionsMap), nil
}
