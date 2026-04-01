package domain

type Section map[string]string

type Config struct {
	Sections map[string]Section
}

func NewConfigFromMap(sections map[string]Section) *Config {
	return &Config{
		Sections: sections,
	}
}

func (c *Config) Get(section, key string) (string, bool) {
	sec, ok := c.Sections[section]
	if !ok {
		return "", false
	}
	value, ok := sec[key]
	return value, ok
}

func (c *Config) Set(section, key, value string) {
	if _, ok := c.Sections[section]; !ok {
		c.Sections[section] = make(Section)
	}
	c.Sections[section][key] = value
}
