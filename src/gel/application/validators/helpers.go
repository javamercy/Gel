package validators

import (
	"regexp"
	"strings"
)

var RegexSHA256 = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)

func isStringSliceNonEmpty(value any) bool {
	paths, ok := value.([]string)
	return ok && len(paths) > 0
}

func areAllInStringSliceNonEmpty(value any) bool {
	paths, ok := value.([]string)
	if !ok {
		return false
	}
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			return false
		}
	}
	return true
}
