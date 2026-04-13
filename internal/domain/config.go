package domain

// Section maps keys to values within one TOML section.
type Section map[string]string

// Config stores repository configuration grouped by TOML section.
type Config struct {
	// Sections maps section names to their key/value pairs.
	Sections map[string]Section
}

// NewConfigFromMap returns a Config backed by sections.
func NewConfigFromMap(sections map[string]Section) *Config {
	return &Config{
		Sections: sections,
	}
}

// Get returns the value for section.key and whether it exists.
func (c *Config) Get(section, key string) (string, bool) {
	sec, ok := c.Sections[section]
	if !ok {
		return "", false
	}

	value, ok := sec[key]
	return value, ok
}

// Set stores value at section.key and creates the section map when needed.
func (c *Config) Set(section, key, value string) {
	if c.Sections == nil {
		c.Sections = make(map[string]Section)
	}
	if _, ok := c.Sections[section]; !ok {
		c.Sections[section] = make(Section)
	}
	c.Sections[section][key] = value
}
