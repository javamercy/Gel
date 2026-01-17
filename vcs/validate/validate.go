package validate

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidHash = errors.New("invalid hash: must be 64 hex characters")
)

var sha256HexRegex = regexp.MustCompile(`^[a-f0-9]{64}$`)

func Hash(hash string) error {
	if !sha256HexRegex.MatchString(hash) {
		return ErrInvalidHash
	}
	return nil
}
