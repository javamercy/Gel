package domain

import (
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	// ErrInvalidHashLength is returned when a hash input has an invalid length.
	ErrInvalidHashLength = errors.New("invalid hash length")

	// ErrInvalidHashEncoding is returned when a hex-encoded hash contains invalid characters.
	ErrInvalidHashEncoding = errors.New("invalid hash encoding")
)

// Hash represents a SHA-256 hash stored as a fixed-length byte array.
type Hash [SHA256ByteLength]byte

// NewHashFromHex creates a new Hash from the provided hexadecimal string and returns an error if the input is invalid.
func NewHashFromHex(hexHash string) (Hash, error) {
	if len(hexHash) != SHA256HexLength {
		return Hash{}, fmt.Errorf(
			"%w: got %d characters, want %d",
			ErrInvalidHashLength,
			len(hexHash),
			SHA256HexLength,
		)
	}

	decoded, err := hex.DecodeString(hexHash)
	if err != nil {
		return Hash{}, fmt.Errorf("%w: %q", ErrInvalidHashEncoding, hexHash)
	}
	return NewHashFromBytes(decoded)
}

// NewHashFromBytes creates a new Hash from the provided byte slice and returns an error if the input has an invalid length.
func NewHashFromBytes(data []byte) (Hash, error) {
	if len(data) != SHA256ByteLength {
		return Hash{}, fmt.Errorf(
			"%w: got %d bytes, want %d",
			ErrInvalidHashLength,
			len(data),
			SHA256ByteLength,
		)
	}

	var hash Hash
	copy(hash[:], data)
	return hash, nil
}

// Hex returns the hex-encoded form of h.
func (h Hash) Hex() string {
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
	return h.Hex()
}
