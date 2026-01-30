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

// FindNullByteIndex finds the index of the first null byte in data.
// Returns -1 if no null byte is found.
func FindNullByteIndex(data []byte) int {
	for i, b := range data {
		if b == 0 {
			return i
		}
	}
	return -1
}
