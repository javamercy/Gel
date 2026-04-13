package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIdentity_Valid(t *testing.T) {
	id, err := NewIdentity("Emre Kurşun", "emre@gmail.com", "1700000000", "+0000")
	require.NoError(t, err)
	assert.Equal(t, "Emre Kurşun", id.Name)
	assert.Equal(t, "emre@gmail.com", id.Email)
	assert.Equal(t, "1700000000", id.Timestamp)
	assert.Equal(t, "+0000", id.Timezone)
}

func TestNewIdentity_Invalid(t *testing.T) {
	_, err := NewIdentity("", "emre@gmail.com", "1700000000", "+0000")
	assert.ErrorIs(t, err, ErrInvalidIdentityFormat)

	_, err = NewIdentity("Emre", "", "1700000000", "+0000")
	assert.ErrorIs(t, err, ErrInvalidIdentityFormat)

	_, err = NewIdentity("Emre", "emre@gmail.com", "not-a-number", "+0000")
	assert.ErrorIs(t, err, ErrInvalidIdentityFormat)

	_, err = NewIdentity("Emre", "emre@gmail.com", "1700000000", "invalid")
	assert.ErrorIs(t, err, ErrInvalidIdentityFormat)
}

func TestNewIdentity_TrimsWhitespace(t *testing.T) {
	id, err := NewIdentity("  Emre  ", "  emre@gmail.com  ", "1700000000", "+0000")
	require.NoError(t, err)
	assert.Equal(t, "Emre", id.Name)
	assert.Equal(t, "emre@gmail.com", id.Email)
}

func TestNewIdentity_TimezoneBounds(t *testing.T) {
	_, err := NewIdentity("Emre", "e@x.com", "1", "+2359")
	require.NoError(t, err)

	_, err = NewIdentity("Emre", "e@x.com", "1", "+2400")
	assert.ErrorIs(t, err, ErrInvalidIdentityFormat)

	_, err = NewIdentity("Emre", "e@x.com", "1", "+2360")
	assert.ErrorIs(t, err, ErrInvalidIdentityFormat)
}
