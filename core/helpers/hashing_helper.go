package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeHash computes the SHA-256 hash of the given data and returns it as a hexadecimal string.
// It is used for generating object hashes in the Gel version control system, following Git's object hashing process.
func ComputeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
