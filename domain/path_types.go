package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

// NormalizedPath represents a path relative to repository root with forward slashes.
type NormalizedPath string

func NewNormalizedPath(repoDir string, path string) (NormalizedPath, error) {
	absPath, err := NewAbsolutePath(path)
	if err != nil {
		return "", fmt.Errorf("failed to create absolute path: %w", err)
	}
	return absPath.ToNormalizedPath(repoDir)
}

// NewNormalizedPathUnchecked creates a NormalizedPath from an already-normalized string
// (e.g., from index file or other storage). It validates the format but doesn't convert.
// Use this when you know the path is already in normalized format.
func NewNormalizedPathUnchecked(path string) (NormalizedPath, error) {
	if err := validateNormalizedFormat(path); err != nil {
		return "", err
	}
	return NormalizedPath(path), nil
}

func (p NormalizedPath) ToAbsolutePath() (AbsolutePath, error) {
	return NewAbsolutePath(p.String())
}

func (p NormalizedPath) String() string {
	return string(p)
}

// validateNormalizedFormat ensures a path is in proper normalized format.
// This provides defense against corrupted index files or malicious data.
func validateNormalizedFormat(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("normalized path cannot be absolute: %s", path)
	}
	if strings.Contains(path, "\\") {
		return fmt.Errorf("normalized path must use forward slashes: %s", path)
	}
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("normalized path cannot contain null bytes: %s", path)
	}
	for i, r := range path {
		if r < 0x20 && r != '\t' {
			return fmt.Errorf("normalized path contains control character at position %d: %s", i, path)
		}
	}
	components := strings.Split(path, "/")
	for i, comp := range components {
		if comp == "" {
			return fmt.Errorf("normalized path cannot contain empty components: %s", path)
		}
		if comp == "." || comp == ".." {
			return fmt.Errorf("normalized path cannot contain %s: %s", comp, path)
		}
		if i == 0 && comp == GelDirName {
			return fmt.Errorf("normalized path cannot start with %s: %s", GelDirName, path)
		}
	}
	return nil
}

// AbsolutePath represents an absolute path to a file or directory.
type AbsolutePath string

func NewAbsolutePath(path string) (AbsolutePath, error) {
	absPath, err := filepath.Abs(filepath.FromSlash(path))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return AbsolutePath(absPath), nil
}

func (p AbsolutePath) ToNormalizedPath(repoDir string) (NormalizedPath, error) {
	relPath, err := filepath.Rel(repoDir, p.String())
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}
	return NormalizedPath(filepath.ToSlash(relPath)), nil
}

func (p AbsolutePath) String() string {
	return string(p)
}
