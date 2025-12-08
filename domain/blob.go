package domain

type Blob struct {
	*BaseObject
}

func NewBlob(data []byte) *Blob {
	return &Blob{
		&BaseObject{
			objectType: GelBlobObjectType,
			data:       data,
		},
	}
}
