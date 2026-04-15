package core

import (
	"Gel/internal/domain"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Exists reports whether path exists on disk.
// It returns (false, nil) when the path is missing and wraps other stat errors.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat path: %w", err)
	}
	return true, nil
}

// SortedPaths sorts a slice of paths (NormalizedPath or AbsolutePath) in lexicographical order based on their string representation.
func SortedPaths[T interface {
	domain.NormalizedPath | domain.AbsolutePath
	String() string
}](paths []T) {
	slices.SortFunc(
		paths, func(a, b T) int {
			return strings.Compare(a.String(), b.String())
		},
	)
}
