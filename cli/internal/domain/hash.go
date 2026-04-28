package domain

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// ErrInvalidHashLen indicates that the provided hash length does not match the expected SHA-256 hash length.
var ErrInvalidHashLen = errors.New("invalid hash length")

// Hash is a fixed-length SHA-256 digest.
type Hash [SHA256ByteLength]byte

// NewHash parses a hex-encoded SHA-256 digest.
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

// ToHexString returns the hex-encoded form of h.
func (h Hash) ToHexString() string {
	return hex.EncodeToString(h[:])
}

// IsEmpty reports whether h is the zero Hash.
func (h Hash) IsEmpty() bool {
	return h == Hash{}
}

// Equals reports whether h and o are identical.
func (h Hash) Equals(o Hash) bool {
	return h == o
}

// String returns the hex-encoded form of h.
func (h Hash) String() string {
	return h.ToHexString()
}
