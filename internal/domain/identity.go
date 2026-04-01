package domain

import (
	"fmt"
)

type Identity struct {
	Name      string
	Email     string
	Timestamp string
	Timezone  string
}

func NewIdentity(name, email, timestamp, timezone string) Identity {
	return Identity{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
		Timezone:  timezone,
	}
}

func (identity Identity) serialize() []byte {
	return []byte(fmt.Sprintf(
		"%s <%s> %s %s",
		identity.Name,
		identity.Email,
		identity.Timestamp,
		identity.Timezone))
}
