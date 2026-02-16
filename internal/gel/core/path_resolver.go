package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

type PathspecType int

func (t PathspecType) String() string {
	switch t {
	case File:
		return "file"
	case Directory:
		return "directory"
	case GlobPattern:
		return "glob pattern"
	case NonExistent:
		return "non-existent"
	}
	return "unknown"
}

var (
	ErrUnknownPathspecType = errors.New("unknown pathspec type")
)

type ResolvedPath struct {
	Type            PathspecType
	NormalizedScope string
	NormalizedPaths map[string]bool
}

type PathResolver struct {
	repositoryDir   string
	ignoredPatterns map[string]bool
}

func NewPathResolver(repositoryDir string, ignoredPatterns map[string]bool) *PathResolver {
	if ignoredPatterns == nil {
		ignoredPatterns = map[string]bool{
			".gel":  true,
			".git":  true,
			".idea": true,
		}
	}
	return &PathResolver{
		repositoryDir:   repositoryDir,
		ignoredPatterns: ignoredPatterns,
	}
}

func (p *PathResolver) Resolve(pathspecs []string) ([]ResolvedPath, error) {
	resolvedPaths := make([]ResolvedPath, 0)

	for _, pathspec := range pathspecs {
		var paths []string
		var err error

		pathspecType := classifyPathspec(pathspec)
		normalizedScope, err := p.normalizePath(pathspec)
		if normalizedScope == "." {
			normalizedScope = ""
		}
		if err != nil {
			return nil, err
		}

		switch pathspecType {
		case File:
			paths = []string{pathspec}
		case Directory:
			paths, err = expandDirectory(pathspec)
		case GlobPattern:
			paths, err = expandGlobPattern(pathspec)
		case NonExistent:
			paths = []string{}
		default:
			return nil, ErrUnknownPathspecType
		}

		if err != nil {
			return nil, err
		}

		normalizedPaths := make(map[string]bool)
		for _, path := range paths {
			normalizedPath, err := p.normalizePath(path)
			if err != nil {
				return nil, err
			}
			// TODO: implement gelignore.
			if p.shouldIgnore(normalizedPath) || normalizedPaths[normalizedPath] {
				continue
			}
			normalizedPaths[normalizedPath] = true
		}

		resolvedPaths = append(
			resolvedPaths, ResolvedPath{
				Type:            pathspecType,
				NormalizedScope: normalizedScope,
				NormalizedPaths: normalizedPaths,
			},
		)
	}

	return resolvedPaths, nil
}

func (p *PathResolver) normalizePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	relPath, err := filepath.Rel(p.repositoryDir, absPath)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(relPath), nil
}

func (p *PathResolver) shouldIgnore(path string) bool {
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if p.ignoredPatterns[segment] {
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

	err := filepath.WalkDir(
		path, func(p string, d os.DirEntry, err error) error {
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
		},
	)
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
