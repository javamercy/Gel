package objects

import "Gel/core/constants"

type Blob struct {
	objectType constants.ObjectType
	size       int
	data       []byte
}

func NewBlob(data []byte) *Blob {
	return &Blob{
		objectType: constants.Blob,
		size:       len(data),
		data:       data,
	}
}

func (b *Blob) GetType() constants.ObjectType {
	return b.objectType
}

func (b *Blob) GetSize() int {
	return b.size
}

func (b *Blob) GetData() []byte {
	return b.data
}
