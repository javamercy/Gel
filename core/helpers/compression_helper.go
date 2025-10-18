package helpers

import (
	"bytes"
	"compress/zlib"
	"io"
)

type ICompressionHelper interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

type ZlibCompressionHelper struct {
}

func NewZlibCompressionHelper() *ZlibCompressionHelper {
	return &ZlibCompressionHelper{}
}

func (zlibCompressionHelper *ZlibCompressionHelper) Decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (zlibCompressionHelper *ZlibCompressionHelper) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		writer.Close()
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
