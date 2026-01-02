package encoding

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

var (
	ErrFailedToDecodeTomlFile = errors.New("failed to decode TOML file")
	ErrFailedToEncodeTomlFile = errors.New("failed to encode TOML file")
)

type ITomlHelper interface {
	Decode(data []byte, value any) error
	Encode(value any) ([]byte, error)
}

type BurntSushiTomlHelper struct {
}

func NewBurntSushiTomlHelper() *BurntSushiTomlHelper {
	return &BurntSushiTomlHelper{}
}

func (burntSushiToml *BurntSushiTomlHelper) Decode(data []byte, value any) error {

	_, err := toml.Decode(string(data), value)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToDecodeTomlFile, err.Error())
	}
	return nil
}

func (burntSushiToml *BurntSushiTomlHelper) Encode(value any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(value); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToEncodeTomlFile, err.Error())
	}
	return buffer.Bytes(), nil
}
