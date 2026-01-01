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

type IToml interface {
	Decode(data []byte, value any) error
	Encode(value any) ([]byte, error)
}

type BurntSushiToml struct {
}

func NewBurntSushiToml() *BurntSushiToml {
	return &BurntSushiToml{}
}

func (burntSushiToml *BurntSushiToml) Decode(data []byte, value any) error {

	_, err := toml.Decode(string(data), value)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToDecodeTomlFile, err.Error())
	}
	return nil
}

func (burntSushiToml *BurntSushiToml) Encode(value any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(value); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToEncodeTomlFile, err.Error())
	}
	return buffer.Bytes(), nil
}
