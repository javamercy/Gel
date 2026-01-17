package validation

import "regexp"

var sha256HexRegex = regexp.MustCompile(`^[a-f0-9]{64}$`)

func IsValidSha256Hex(hash string) bool {
	return sha256HexRegex.MatchString(hash)
}
