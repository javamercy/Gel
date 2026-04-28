package domain

// Blob represents immutable object content.
type Blob struct {
	body []byte
}

// Body returns a copy of the blob contents.
func (blob *Blob) Body() []byte {
	return append([]byte(nil), blob.body...)
}

// NewBlob returns a Blob containing a copy of body.
func NewBlob(body []byte) *Blob {
	blob := &Blob{
		body: append([]byte(nil), body...),
	}
	return blob
}

// Type returns ObjectTypeBlob.
func (blob *Blob) Type() ObjectType {
	return ObjectTypeBlob
}

// Size returns the size of the blob contents in bytes.
func (blob *Blob) Size() int {
	return len(blob.body)
}

// Serialize returns the full object serialization in the form "<type> <size>\x00<body>".
func (blob *Blob) Serialize() []byte {
	return SerializeObject(ObjectTypeBlob, blob.body)
}
