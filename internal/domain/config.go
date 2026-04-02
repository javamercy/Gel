package domain

// Section is a key-value map for one TOML section.
type Section map[string]string

// Config stores all repository configuration sections.
type Config struct {
	Sections map[string]Section
}

// NewConfigFromMap constructs Config from decoded section map data.
func NewConfigFromMap(sections map[string]Section) *Config {
	return &Config{
		Sections: sections,
	}
}

// Get returns section.key and whether it exists.
func (c *Config) Get(section, key string) (string, bool) {
	sec, ok := c.Sections[section]
	if !ok {
		return "", false
	}

	value, ok := sec[key]
	return value, ok
}

// Set assigns section.key=value, creating the section when absent.
func (c *Config) Set(section, key, value string) {
	if _, ok := c.Sections[section]; !ok {
		c.Sections[section] = make(Section)
	}
	c.Sections[section][key] = value
}
