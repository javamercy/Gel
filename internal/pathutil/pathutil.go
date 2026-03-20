package pathutil

import (
	"errors"
	"fmt"
	"os"
)

func IsFile(path string) (bool, error) {
	file, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat path: %w", err)
	}
	if file.IsDir() {
		return false, nil
	}
	return true, nil
}

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
