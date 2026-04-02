package core

import (
	"Gel/internal/domain"
	"Gel/internal/storage"
	"Gel/internal/validate"
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/BurntSushi/toml"
)

const (
	// ConfigSectionUser stores author/committer identity defaults.
	ConfigSectionUser = "user"
	// ConfigKeyName is the user name key under [user].
	ConfigKeyName = "name"
	// ConfigKeyEmail is the user email key under [user].
	ConfigKeyEmail = "email"
)

// ConfigService manages repository config stored in .gel/config.toml.
type ConfigService struct {
	configStorage *storage.ConfigStorage
}

// NewConfigService creates a config service.
func NewConfigService(configStorage *storage.ConfigStorage) *ConfigService {
	return &ConfigService{
		configStorage: configStorage,
	}
}

// GetUserInfo returns user.name and user.email used by commit operations.
func (c *ConfigService) GetUserInfo() (string, string, error) {
	config, err := c.Read()
	if err != nil {
		return "", "", err
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

// Set writes section.key=value to config, creating missing sections as needed.
func (c *ConfigService) Set(section, key, value string) error {
	if err := validate.StringMustNotBeEmpty(section); err != nil {
		return fmt.Errorf("config: %w", err)
	}
	if err := validate.StringMustNotBeEmpty(key); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	config, err := c.Read()
	if err != nil {
		return err
	}
	config.Set(section, key, value)
	return c.Write(config)
}

// GetAndOutput reads a config value and prints it to writer.
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

// Get returns a config value by section and key.
func (c *ConfigService) Get(section, key string) (string, error) {
	if err := validate.StringMustNotBeEmpty(section); err != nil {
		return "", fmt.Errorf("config: %w", err)
	}
	if err := validate.StringMustNotBeEmpty(key); err != nil {
		return "", fmt.Errorf("config: %w", err)
	}

	config, err := c.Read()
	if err != nil {
		return "", err
	}
	value, ok := config.Get(section, key)
	if !ok {
		return "", fmt.Errorf("config: key '%s.%s' not found", section, key)
	}
	return value, nil
}

// List writes all config entries in "section.key=value" format.
// Output is sorted by section then key for deterministic ordering.
func (c *ConfigService) List(writer io.Writer) error {
	config, err := c.Read()
	if err != nil {
		return err
	}

	sectionNames := make([]string, 0, len(config.Sections))
	for sectionName := range config.Sections {
		sectionNames = append(sectionNames, sectionName)
	}
	sort.Strings(sectionNames)

	for _, sectionName := range sectionNames {
		section := config.Sections[sectionName]

		keys := make([]string, 0, len(section))
		for key := range section {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := section[key]
			if _, err := fmt.Fprintf(writer, "%s.%s=%s\n", sectionName, key, value); err != nil {
				return fmt.Errorf("config: failed to write config: %w", err)
			}
		}
	}
	return nil
}

// Write encodes and persists config data to storage.
func (c *ConfigService) Write(config *domain.Config) error {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(config.Sections); err != nil {
		return fmt.Errorf("config: failed to encode config: %w", err)
	}
	if err := c.configStorage.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// Read loads and decodes config from storage.
// Empty files are treated as empty config maps.
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
		return nil, fmt.Errorf("config: failed to decode config: %w", err)
	}
	return domain.NewConfigFromMap(sectionsMap), nil
}
