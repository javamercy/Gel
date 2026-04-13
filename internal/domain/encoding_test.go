package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeSHA256(t *testing.T) {
	data := []byte("hello world")
	sum := sha256.Sum256(data)
	expected := hex.EncodeToString(sum[:])

	assert.Equal(t, expected, ComputeSHA256(data))
}
