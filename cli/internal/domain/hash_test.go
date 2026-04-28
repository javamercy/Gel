package domain

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHash(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr error
		wantErr   bool
	}{
		{name: "valid", input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
		{name: "invalid length", input: "abcd", expectErr: ErrInvalidHashLen},
		{name: "invalid hex", input: "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				h, err := NewHash(tt.input)
				if tt.expectErr != nil || tt.wantErr {
					assert.Error(t, err)
					if tt.expectErr == ErrInvalidHashLen {
						assert.ErrorIs(t, err, ErrInvalidHashLen)
					}
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tt.input, h.ToHexString())
			},
		)
	}
}

func TestHashHelpers(t *testing.T) {
	h, err := NewHash("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	require.NoError(t, err)

	assert.False(t, h.IsEmpty())
	assert.True(t, h.Equals(h))
	assert.Equal(t, h.ToHexString(), h.String())

	var zero Hash
	assert.True(t, zero.IsEmpty())
	assert.False(t, zero.Equals(h))
}

func TestHashToHexString_Compatibility(t *testing.T) {
	raw := make([]byte, SHA256ByteLength)
	for i := range raw {
		raw[i] = byte(i)
	}

	var h Hash
	copy(h[:], raw)
	assert.Equal(t, hex.EncodeToString(raw), h.ToHexString())
}
