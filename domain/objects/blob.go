package objects

import "Gel/core/constants"

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
