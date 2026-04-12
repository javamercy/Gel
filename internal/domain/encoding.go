package domain

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeSHA256 computes the SHA-256 hash of data and returns it as hex string.
// This is placed in domain as it's used for computing checksums of domain objects.
func ComputeSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
