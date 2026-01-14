package util

import (
	"Gel/core/constant"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type PathspecType int

var (
	ErrPathspecDidNotMatchAny = errors.New("pathspec did not match any file, directory, or glob pattern")
	ErrUnknownPathspecType    = errors.New("unknown pathspec type")
)

const (
	File PathspecType = iota
	Directory
	GlobPattern
	NonExistent
)

const (
	globPatterns string = "*?[]"
)

type IPathResolver interface {
	Resolve(pathspecs []string) ([]string, error)
}

var _ IPathResolver = (*PathResolver)(nil)

type PathResolver struct {
	repositoryDirectory string
	ignoredPatterns     map[string]bool
}

func NewPathResolver(repositoryDirectory string, ignoredPatterns map[string]bool) *PathResolver {
	if ignoredPatterns == nil {
		ignoredPatterns = map[string]bool{
			".gel":  true,
			".git":  true,
			".idea": true,
		}
	}
	return &PathResolver{
		repositoryDirectory: repositoryDirectory,
		ignoredPatterns:     ignoredPatterns,
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

			// TODO: implement gelignore.
			if pathResolver.shouldIgnore(normalizedPath) || normalizedPathMap[normalizedPath] {
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
		return nil, ErrPathspecDidNotMatchAny
	default:
		return nil, ErrUnknownPathspecType
	}
}

func (pathResolver *PathResolver) normalizePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	relPath, err := filepath.Rel(pathResolver.repositoryDirectory, absPath)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(relPath), nil
}

func (pathResolver *PathResolver) shouldIgnore(path string) bool {
	segments := strings.Split(path, constant.SlashStr)

	for _, segment := range segments {
		if pathResolver.ignoredPatterns[segment] {
			return true
		}
	}
	return false
}

func classifyPathspec(pathspec string) PathspecType {
	if strings.ContainsAny(pathspec, globPatterns) {
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

func expandDirectory(path string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := os.Stat(p)
			if err == nil && info.IsDir() {
				return nil
			}
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func expandGlobPattern(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if info.IsDir() {
			files, err := expandDirectory(match)
			if err != nil {
				return nil, err
			}
			result = append(result, files...)
		} else {
			result = append(result, match)
		}
	}
	return result, nil
}
