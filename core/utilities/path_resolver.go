package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PathspecType int

const (
	File PathspecType = iota
	Directory
	GlobPattern
	NonExistent
)

type PathResolver struct {
	repositoryDir string
}

func NewPathResolver(repositoryDir string) *PathResolver {
	return &PathResolver{
		repositoryDir: repositoryDir,
	}
}

func (pathResolver *PathResolver) Resolve(pathspecs []string) ([]string, error) {
	normalizedPathMap := make(map[string]bool)
	var normalizedPaths []string

	for _, pathspec := range pathspecs {
		paths, err := pathResolver.resolvePathspec(pathspec)
		if err != nil {
			return nil, err
		}

		for _, path := range paths {

			normalizedPath, err := pathResolver.normalizePath(path)
			if err != nil {
				return nil, err
			}

			// TODO: Bypass .gel, .git, and .idea directories for now. Implement .gelignore for later.
			if strings.Contains(path, ".gel"+string(os.PathSeparator)) ||
				strings.Contains(path, ".git"+string(os.PathSeparator)) ||
				strings.Contains(path, ".idea"+string(os.PathSeparator)) {
				continue
			}

			if normalizedPathMap[normalizedPath] {
				continue
			}
			normalizedPathMap[normalizedPath] = true

			normalizedPaths = append(normalizedPaths, normalizedPath)
		}
	}

	return normalizedPaths, nil
}

func (pathResolver *PathResolver) resolvePathspec(pathspec string) ([]string, error) {

	switch classifyPathspec(pathspec) {
	case File:
		return []string{pathspec}, nil
	case Directory:
		return expandDirectory(pathspec)
	case GlobPattern:
		return expandGlobPattern(pathspec)
	case NonExistent:
		return nil, fmt.Errorf("pathspec '%s' did not match any files", pathspec)
	default:
		return nil, fmt.Errorf("unknown pathspec type for '%s'", pathspec)
	}
}

func (pathResolver *PathResolver) normalizePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	relPath, err := filepath.Rel(pathResolver.repositoryDir, absPath)
	if err != nil {
		return "", err
	}

	return relPath, nil
}

func expandDirectory(path string) ([]string, error) {
	var files []string

	walkDirErr := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			files = append(files, p)
		}
		return nil
	})

	if walkDirErr != nil {
		return nil, walkDirErr
	}

	return files, nil
}

func classifyPathspec(pathspec string) PathspecType {
	var globPatternsString = "*?[]"
	if strings.ContainsAny(pathspec, globPatternsString) {
		return GlobPattern
	}

	fileInfo, err := os.Stat(pathspec)
	if os.IsNotExist(err) {
		return NonExistent
	}

	if fileInfo.IsDir() {
		return Directory
	}
	return File
}

func expandGlobPattern(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
