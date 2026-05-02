package domain

import (
	"slices"
	"strings"
)

type sortablePath interface {
	NormalizedPath | AbsolutePath
	String() string
}

// SortPaths sorts paths in-place by their string representation.
func SortPaths[T sortablePath](paths []T) {
	slices.SortFunc(
		paths, func(a, b T) int {
			return strings.Compare(a.String(), b.String())
		},
	)
}

// SortedPathSet returns the keys from paths sorted by their string representation.
func SortedPathSet[T sortablePath](paths map[T]struct{}) []T {
	out := make([]T, 0, len(paths))
	for path := range paths {
		out = append(out, path)
	}
	SortPaths(out)
	return out
}
