package domain

import (
	"Gel/core/constant"
	"bytes"
)

type Identity struct {
	Name      string
	Email     string
	Timestamp string
	Timezone  string
}

func NewIdentity(name, email, timestamp, timezone string) *Identity {
	return &Identity{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
		Timezone:  timezone,
	}
}

func (identity *Identity) serialize() []byte {
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
