package objects

import (
	"Gel/src/gel/core/constants"
)

type Blob struct {
	*BaseObject
}

func NewBlob(data []byte) *Blob {
	return &Blob{
		&BaseObject{
			objectType: constants.Blob,
			data:       data,
		},
	}
}
