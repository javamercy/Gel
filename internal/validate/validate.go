package validate

import (
	"errors"
	"fmt"
	"os"
)

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

func StringMustNotBeEmpty(s string) error {
	if len(s) == 0 {
		return errors.New("string must not be empty")
	}
	return nil
}
