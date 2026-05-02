package domain

// Blob represents immutable file content stored as a blob object.
type Blob struct {
	body []byte
}

// NewBlob returns a Blob containing a defensive copy of body.
func NewBlob(body []byte) *Blob {
	blob := &Blob{
		body: append([]byte(nil), body...),
	}
	return blob
}

// Body returns a defensive copy of the blob contents.
func (blob *Blob) Body() []byte {
	return append([]byte(nil), blob.body...)
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
