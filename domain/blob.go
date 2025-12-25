package domain

import "Gel/core/validation"

type Blob struct {
	body []byte `validate:"required"`
}

func (blob *Blob) Body() []byte {
	return blob.body
}

func NewBlob(body []byte) (*Blob, error) {
	validator := validation.GetValidator()
	blob := &Blob{
		body: body,
	}
	if err := validator.Struct(blob); err != nil {
		return nil, err
	}
	return blob, nil
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
