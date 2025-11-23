package objects

import (
	"Gel/src/gel/core/constant"
)

type Blob struct {
	*BaseObject
}

func NewBlob(data []byte) *Blob {
	return &Blob{
		&BaseObject{
			objectType: constant.GelBlobObjectType,
			data:       data,
		},
	}
}
