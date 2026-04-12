package domain

type Blob struct {
	body []byte
}

func (blob *Blob) Body() []byte {
	return append([]byte(nil), blob.body...)
}

func NewBlob(body []byte) *Blob {
	blob := &Blob{
		body: append([]byte(nil), body...),
	}
	return blob
}

func (blob *Blob) Type() ObjectType {
	return ObjectTypeBlob
}

func (blob *Blob) Size() int {
	return len(blob.body)
}

func (blob *Blob) Serialize() []byte {
	return SerializeObject(ObjectTypeBlob, blob.body)
}
