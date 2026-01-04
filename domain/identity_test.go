package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserIdentity_Valid(t *testing.T) {
	id, err := NewUserIdentity("Emre Kurşun", "emre@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, "Emre Kurşun", id.Name)
	assert.Equal(t, "emre@gmail.com", id.Email)
}

func TestNewUserIdentity_Invalid(t *testing.T) {
	_, err := NewUserIdentity("", "emre@gmail.com")
	assert.Error(t, err)

	_, err = NewUserIdentity("Emre Kurşun", "invalid-email")
	assert.Error(t, err)
}

func TestNewIdentity_Valid(t *testing.T) {
	id, err := NewIdentity("Emre Kurşun", "emre@gmail.com", "1700000000", "+0000")
	assert.NoError(t, err)
	assert.Equal(t, "Emre Kurşun", id.User.Name)
	assert.Equal(t, "emre@gmail.com", id.User.Email)
	assert.Equal(t, "1700000000", id.Timestamp)
	assert.Equal(t, "+0000", id.Timezone)
}

func TestNewIdentity_Invalid(t *testing.T) {
	_, err := NewIdentity("", "emre@gmail.com", "1700000000", "+0000")
	assert.Error(t, err)

	_, err = NewIdentity("Emre Kurşun", "emre@gmail.com", "", "+0000")
	assert.Error(t, err)

	_, err = NewIdentity("Emre Kurşun", "emre@gmail.com", "1700000000", "")
	assert.Error(t, err)
}
