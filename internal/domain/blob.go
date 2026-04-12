package domain

// Blob represents a data structure encapsulating an immutable byte slice as its content.
type Blob struct {
	body []byte
}

// Body returns a copy of the blob's underlying byte slice to ensure immutability of the original data.
func (blob *Blob) Body() []byte {
	return append([]byte(nil), blob.body...)
}

// NewBlob creates a new Blob object with a copy of the provided byte slice to ensure immutability of the original data.
func NewBlob(body []byte) *Blob {
	blob := &Blob{
		body: append([]byte(nil), body...),
	}
	return blob
}

// Type returns the object type of the Blob, which is always ObjectTypeBlob.
func (blob *Blob) Type() ObjectType {
	return ObjectTypeBlob
}

// Size returns the size of the blob's body in bytes.
func (blob *Blob) Size() int {
	return len(blob.body)
}

// Serialize converts the Blob object into a serialized byte slice in the format "<type> <size>\x00<body>".
func (blob *Blob) Serialize() []byte {
	return SerializeObject(ObjectTypeBlob, blob.body)
}
