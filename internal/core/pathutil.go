package core

import (
	"errors"
	"fmt"
	"os"
)

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
