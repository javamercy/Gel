package domain

import (
	"slices"
	"strings"
)

// SortedPaths sorts a slice of paths (NormalizedPath or AbsolutePath) in lexicographical order based on their string representation.
func SortedPaths[T interface {
	NormalizedPath | AbsolutePath
	String() string
}](paths []T) {
	slices.SortFunc(
		paths, func(a, b T) int {
			return strings.Compare(a.String(), b.String())
		},
	)
}

func SortedPathSet[T interface {
	NormalizedPath | AbsolutePath
	String() string
}](paths map[T]bool) []T {
	out := make([]T, 0, len(paths))
	for path := range paths {
		out = append(out, path)
	}
	SortedPaths(out)
	return out
}
