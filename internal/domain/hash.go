package domain

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// ErrInvalidHashLen indicates that the provided hash length does not match the expected SHA-256 hash length.
var ErrInvalidHashLen = errors.New("invalid hash length")

// Hash represents a fixed-length SHA-256 hash as a 32-byte array.
type Hash [SHA256ByteLength]byte

// NewHash parses a SHA-256 hash string and returns a Hash object or an error if the string is invalid.
func NewHash(s string) (Hash, error) {
	if len(s) != SHA256HexLength {
		return Hash{}, ErrInvalidHashLen
	}

	decoded, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, fmt.Errorf("failed to decode hash: %w", err)
	}
	return Hash(decoded), nil
}

// ToHexString converts the Hash to its hexadecimal string representation.
func (h Hash) ToHexString() string {
	return hex.EncodeToString(h[:])
}

// IsEmpty returns true if the Hash is uninitialized or empty (i.e., equals the zero-value Hash).
func (h Hash) IsEmpty() bool {
	return h == Hash{}
}

// Equals checks if the current Hash is equal to another Hash. Returns true if both are identical, otherwise false.
func (h Hash) Equals(o Hash) bool {
	return h == o
}

// String returns the hexadecimal string representation of the Hash.
func (h Hash) String() string {
	return h.ToHexString()
}
