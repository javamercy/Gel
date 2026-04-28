package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlob_Valid(t *testing.T) {
	body := []byte("hello world")
	blob := NewBlob(body)
	assert.NotNil(t, blob)
	assert.Equal(t, body, blob.Body())
	assert.Equal(t, ObjectTypeBlob, blob.Type())
	assert.Equal(t, len(body), blob.Size())
}

func TestNewBlob_NilBody(t *testing.T) {
	blob := NewBlob(nil)
	assert.NotNil(t, blob)
	assert.Equal(t, 0, blob.Size())
}

func TestBlob_Immutability(t *testing.T) {
	input := []byte("hello")
	blob := NewBlob(input)

	input[0] = 'H'
	assert.Equal(t, []byte("hello"), blob.Body())

	body := blob.Body()
	body[1] = 'A'
	assert.Equal(t, []byte("hello"), blob.Body())
}
