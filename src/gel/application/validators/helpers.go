package validators

import (
	"regexp"
	"strings"
)

var regexSHA256 = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)

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

func exactlyOne(values ...bool) bool {
	count := 0
	for _, value := range values {
		if value {
			count++
		}
	}
	return count == 1
}

func atLeastOne(values ...bool) bool {
	for _, value := range values {
		if value {
			return true
		}
	}
	return false
}
