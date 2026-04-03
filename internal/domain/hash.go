package domain

import (
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	ErrInvalidHashLen = errors.New("invalid hash length")
)

type Hash [SHA256ByteLength]byte

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

func (h Hash) ToHexString() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) IsEmpty() bool {
	return h == Hash{}
}

func (h Hash) Equals(o Hash) bool {
	return h == o
}

func (h Hash) String() string {
	return h.ToHexString()
}
