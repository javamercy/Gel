package domain

import (
	"Gel/core/constant"
	"Gel/core/validation"
	"bytes"
)

type Identity struct {
	Name      string `validate:"required,min=1,max=256"`
	Email     string `validate:"required,email"`
	Timestamp string `validate:"required"`
	Timezone  string `validate:"required,timezone"`
}

func NewIdentity(name, email, timestamp, timezone string) (Identity, error) {
	identity := Identity{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
		Timezone:  timezone,
	}

	validator := validation.GetValidator()
	if err := validator.Struct(identity); err != nil {
		return Identity{}, err
	}

	return identity, nil
}

func (identity Identity) serialize() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(identity.Name)
	buffer.WriteByte(constant.SpaceByte)
	buffer.WriteByte(constant.LessThanByte)
	buffer.WriteString(identity.Email)
	buffer.WriteByte(constant.GreaterThanByte)
	buffer.WriteByte(constant.SpaceByte)
	buffer.WriteString(identity.Timestamp)
	buffer.WriteByte(constant.SpaceByte)
	buffer.WriteString(identity.Timezone)

	return buffer.Bytes()
}
