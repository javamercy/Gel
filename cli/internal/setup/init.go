package setup

import (
	"Gel/internal/domain"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	// ErrInitEmptyPath is returned when repository initialization receives an
	// empty path.
	ErrInitEmptyPath = errors.New("init path is empty")

	// ErrInitPathNotDirectory is returned when a required initialization
	// directory path exists as a non-directory.
	ErrInitPathNotDirectory = errors.New("path is not a directory")

	// ErrInitPathNotRegularFile is returned when a required initialization file
	// path exists as a non-regular file.
	ErrInitPathNotRegularFile = errors.New("path is not a regular file")
)

// InitService bootstraps the filesystem layout for a Gel repository.
//
// InitService is stateless. It owns only the repository creation workflow; it
// does not discover existing repositories or initialize services that require a
// repository to already exist.
type InitService struct {
}

// NewInitService creates a stateless service for repository initialization.
// The service only performs filesystem operations and keeps no in-memory state.
func NewInitService() *InitService {
	return &InitService{}
}

// Init bootstraps a Gel repository at the provided path.
//
// Behavior:
//   - Validates that path is not empty and resolves it to an absolute path.
//   - Creates the target directory if it does not exist.
//   - Fails when the target exists but is not a directory.
//   - Ensures .gel, .gel/objects, and .gel/refs/heads directories exist.
//   - Creates .gel/HEAD with "ref: refs/heads/main" when missing.
//   - Creates an empty .gel/config.toml when missing.
//   - Leaves existing HEAD and config files unchanged when they are regular
//     files.
//
// Init is idempotent: running it repeatedly on an existing repository is valid
// and returns a reinitialization message.
//
// Returned errors wrap ErrInitEmptyPath, ErrInitPathNotDirectory, or
// ErrInitPathNotRegularFile for expected initialization conflicts. Filesystem
// failures are wrapped with the operation and path that failed.
func (i *InitService) Init(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("init: validate path: %w", ErrInitEmptyPath)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("init: resolve path %q: %w", path, err)
	}

	gelPath := filepath.Join(absPath, domain.GelDirName)
	objectsPath := filepath.Join(gelPath, domain.ObjectsDirName)
	headsPath := filepath.Join(gelPath, domain.RefsDirName, domain.HeadsDirName)
	headPath := filepath.Join(gelPath, domain.HeadFileName)
	configPath := filepath.Join(gelPath, domain.ConfigFileName)

	gelExists, err := directoryExists(gelPath)
	if err != nil {
		return "", fmt.Errorf("init: check existing repository: %w", err)
	}
	if err := ensureDirectory(absPath); err != nil {
		return "", fmt.Errorf("init: prepare repository directory: %w", err)
	}
	if err := ensureDirectory(gelPath); err != nil {
		return "", fmt.Errorf("init: prepare metadata directory: %w", err)
	}
	if err := ensureDirectory(objectsPath); err != nil {
		return "", fmt.Errorf("init: prepare objects directory: %w", err)
	}
	if err := ensureDirectory(headsPath); err != nil {
		return "", fmt.Errorf("init: prepare refs directory: %w", err)
	}

	headRefContent := fmt.Sprintf("ref: %s\n", domain.MainRef)
	if err := ensureRegularFile(headPath, []byte(headRefContent)); err != nil {
		return "", fmt.Errorf("init: prepare HEAD: %w", err)
	}
	if err := ensureRegularFile(configPath, nil); err != nil {
		return "", fmt.Errorf("init: prepare config: %w", err)
	}

	if gelExists {
		return fmt.Sprintf("Reinitialized existing Gel repository in %v", gelPath), nil
	}
	return fmt.Sprintf("Initialized empty Gel repository in %v", gelPath), nil
}

// ensureDirectory creates path when it is missing and verifies that an existing
// path is a directory.
func ensureDirectory(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%q: %w", path, ErrInitPathNotDirectory)
		}
		return nil
	}
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, domain.DirPermission); err != nil {
			return fmt.Errorf("create directory %q: %w", path, err)
		}
		return nil
	}
	return fmt.Errorf("stat %q: %w", path, err)
}

// ensureRegularFile writes body to path when it is missing and verifies that an
// existing path is a regular file. Existing regular files are left unchanged.
func ensureRegularFile(path string, body []byte) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("%q: %w", path, ErrInitPathNotRegularFile)
		}
		return nil
	}
	if errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(path, body, domain.FilePermission); err != nil {
			return fmt.Errorf("create file %q: %w", path, err)
		}
		return nil
	}
	return fmt.Errorf("stat %q: %w", path, err)
}

// directoryExists reports whether path currently exists as a directory.
func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("stat %q: %w", path, err)
}
