package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIdentity_Valid(t *testing.T) {
	id := NewIdentity("Emre Kurşun", "emre@gmail.com", "1700000000", "+0000")
	assert.Equal(t, "Emre Kurşun", id.Name)
	assert.Equal(t, "emre@gmail.com", id.Email)
	assert.Equal(t, "1700000000", id.Timestamp)
	assert.Equal(t, "+0000", id.Timezone)
}
