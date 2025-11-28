package validators

import "strings"

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
