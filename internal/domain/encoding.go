package domain

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeSHA256 calculates the SHA-256 hash of the input data and returns it as a hexadecimal string.
func ComputeSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
