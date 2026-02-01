package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFromMap_Valid(t *testing.T) {
	sections := map[string]Section{
		"user": {
			"name":  "Emre Kurşun",
			"email": "emre@gmail.com",
		},
	}
	config := NewConfigFromMap(sections)
	assert.NotNil(t, config)

	val, ok := config.Get("user", "name")
	assert.True(t, ok)
	assert.Equal(t, "Emre Kurşun", val)

	val, ok = config.Get("user", "email")
	assert.True(t, ok)
	assert.Equal(t, "emre@gmail.com", val)
}

func TestConfig_Set(t *testing.T) {
	config := NewConfigFromMap(make(map[string]Section))

	config.Set("core", "editor", "vim")

	val, ok := config.Get("core", "editor")
	assert.True(t, ok)
	assert.Equal(t, "vim", val)
}

func TestConfig_Set_ExistingSection(t *testing.T) {
	sections := map[string]Section{
		"user": {
			"name": "Old Name",
		},
	}
	config := NewConfigFromMap(sections)

	config.Set("user", "name", "New Name")

	val, ok := config.Get("user", "name")
	assert.True(t, ok)
	assert.Equal(t, "New Name", val)
}
