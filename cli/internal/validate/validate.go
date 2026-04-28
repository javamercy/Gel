package validate

import (
	"errors"
	"fmt"
	"os"
)

// PathMustExist returns an error when path does not exist.
func PathMustExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("path %s does not exist", path)
		}
		return err
	}
	return nil
}

// PathMustBeDirectory returns an error unless path exists and is a directory.
func PathMustBeDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsDir() {
		return fmt.Errorf("path %s is not a directory", path)
	}
	return nil
}

// PathMustBeFile returns an error unless path exists and is a regular file.
func PathMustBeFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("path %s is not a file", path)
	}
	return nil
}

// StringMustNotBeEmpty returns an error when s is empty.
func StringMustNotBeEmpty(s string) error {
	if len(s) == 0 {
		return errors.New("string must not be empty")
	}
	return nil
}
