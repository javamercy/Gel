package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserConfig_Valid(t *testing.T) {
	config, err := NewUserConfig("Emre Kurşun", "emre@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, "Emre Kurşun", config.Name)
	assert.Equal(t, "emre@gmail.com", config.Email)
}

func TestNewUserConfig_InvalidName(t *testing.T) {
	_, err := NewUserConfig("", "emre@gmail.com")
	assert.Error(t, err)
}

func TestNewUserConfig_InvalidEmail(t *testing.T) {
	_, err := NewUserConfig("Emre Kurşun", "invalid-email")
	assert.Error(t, err)
}

func TestNewConfig_Valid(t *testing.T) {
	userConfig, _ := NewUserConfig("Emre Kurşun", "emre@gmail.com")
	config, err := NewConfig(userConfig)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, userConfig, config.User)
}

func TestNewConfig_Invalid(t *testing.T) {
	config, err := NewConfig(UserConfig{})
	assert.Error(t, err)
	assert.Nil(t, config)
}
