package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfigKey is returned when a config section or key is invalid.
	ErrInvalidConfigKey = errors.New("invalid config key")
)

// ConfigSection maps keys to values within one TOML section.
type ConfigSection map[string]string

// ConfigEntry represents a single configuration entry with section, key, and value fields.
type ConfigEntry struct {
	Section string
	Key     string
	Value   string
}

// Config stores repository configuration grouped by TOML section.
type Config struct {
	// sections maps section names to their key/value pairs.
	sections map[string]ConfigSection
}

// NewConfig returns an empty Config.
func NewConfig() *Config {
	return &Config{
		sections: make(map[string]ConfigSection),
	}
}

// NewConfigFromSections returns a Config containing a defensive copy of sections.
func NewConfigFromSections(sections map[string]ConfigSection) (*Config, error) {
	config := NewConfig()
	for sectionName, section := range sections {
		if err := validateConfigName("section", sectionName); err != nil {
			return nil, err
		}
		for key, value := range section {
			if err := validateConfigName("key", key); err != nil {
				return nil, err
			}
			config.setUnchecked(sectionName, key, value)
		}
	}
	return config, nil
}

// Get returns the value for section.key and whether it exists.
func (c *Config) Get(section, key string) (string, bool) {
	if c == nil || c.sections == nil {
		return "", false
	}

	sec, ok := c.sections[section]
	if !ok {
		return "", false
	}

	value, ok := sec[key]
	return value, ok
}

// Set stores value at section.key and creates the section map when needed.
func (c *Config) Set(section, key, value string) error {
	if err := validateConfigName("section", section); err != nil {
		return err
	}
	if err := validateConfigName("key", key); err != nil {
		return err
	}
	if c.sections == nil {
		c.sections = make(map[string]ConfigSection)
	}
	c.setUnchecked(section, key, value)
	return nil
}

// Entries returns all config entries in unspecified order.
func (c *Config) Entries() []ConfigEntry {
	if c == nil || c.sections == nil {
		return nil
	}

	entries := make([]ConfigEntry, 0)
	for sectionName, section := range c.sections {
		for key, value := range section {
			entries = append(
				entries, ConfigEntry{
					Section: sectionName,
					Key:     key,
					Value:   value,
				},
			)
		}
	}
	return entries
}

// Sections returns a defensive copy of all config sections.
func (c *Config) Sections() map[string]ConfigSection {
	if c == nil || c.sections == nil {
		return make(map[string]ConfigSection)
	}
	sections := make(map[string]ConfigSection, len(c.sections))
	for sectionName, section := range c.sections {
		sections[sectionName] = cloneConfigSection(section)
	}
	return sections
}

func (c *Config) setUnchecked(section, key, value string) {
	if _, ok := c.sections[section]; !ok {
		c.sections[section] = make(ConfigSection)
	}
	c.sections[section][key] = value
}

func validateConfigName(kind, value string) error {
	if value == "" {
		return fmt.Errorf("%w: %s is empty", ErrInvalidConfigKey, kind)
	}
	return nil
}

func cloneConfigSection(section ConfigSection) ConfigSection {
	cloned := make(ConfigSection, len(section))
	for key, value := range section {
		cloned[key] = value
	}
	return cloned
}
