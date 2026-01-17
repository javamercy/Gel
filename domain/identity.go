package domain

import (
	"Gel/domain/validation"
	"fmt"
)

type UserIdentity struct {
	Name  string `validate:"required,min=1,max=256"`
	Email string `validate:"required,email"`
}

func NewUserIdentity(name, email string) (UserIdentity, error) {
	userIdentity := UserIdentity{
		Name:  name,
		Email: email,
	}

	validator := validation.GetValidator()
	if err := validator.Struct(userIdentity); err != nil {
		return UserIdentity{}, err
	}

	return userIdentity, nil
}

type Identity struct {
	User      UserIdentity `validate:"required"`
	Timestamp string       `validate:"required"`
	Timezone  string       `validate:"required,timezone"`
}

func NewIdentity(name, email, timestamp, timezone string) (Identity, error) {

	identity := Identity{
		User: UserIdentity{
			Name:  name,
			Email: email,
		},
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
	return []byte(fmt.Sprintf(
		"%s <%s> %s %s",
		identity.User.Name,
		identity.User.Email,
		identity.Timestamp,
		identity.Timezone))
}
