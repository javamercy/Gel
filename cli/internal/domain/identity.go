package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrInvalidIdentity is returned when an identity cannot be serialized safely.
	ErrInvalidIdentity = errors.New("invalid identity")
)

// Identity stores commit author or committer metadata.
type Identity struct {
	// Name is the display name of the actor.
	Name string

	// Email is the email address of the actor.
	Email string

	// Timestamp is a Unix timestamp in seconds, encoded as base-10 string.
	Timestamp string

	// Timezone is an offset in Git-style ±HHMM format (for example, +0300).
	Timezone string
}

// NewIdentity validates and returns an Identity.
// Validation requires non-empty name/email, a minimal email shape, a valid Unix
// timestamp string, and a valid ±HHMM timezone.
func NewIdentity(name, email, timestamp, timezone string) (Identity, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if err := validateIdentityName(name); err != nil {
		return Identity{}, err
	}
	if err := validateIdentityEmail(email); err != nil {
		return Identity{}, err
	}
	if err := validateIdentityTimestamp(timestamp); err != nil {
		return Identity{}, err
	}
	if err := validateIdentityTimezone(timezone); err != nil {
		return Identity{}, err
	}
	return Identity{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
		Timezone:  timezone,
	}, nil
}

// Serialize returns identity in commit-header form:
// "<name> <email> <timestamp> <timezone>".
func (i Identity) Serialize() []byte {
	return []byte(fmt.Sprintf(
		"%s <%s> %s %s",
		i.Name,
		i.Email,
		i.Timestamp,
		i.Timezone,
	))
}

func validateIdentityName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is empty", ErrInvalidIdentity)
	}
	if strings.ContainsAny(name, "<>\n\x00") {
		return fmt.Errorf("%w: name contains reserved characters", ErrInvalidIdentity)
	}
	return nil
}

func validateIdentityEmail(email string) error {
	if email == "" {
		return fmt.Errorf("%w: email is empty", ErrInvalidIdentity)
	}
	if strings.ContainsAny(email, "<> \t\n\x00") {
		return fmt.Errorf("%w: email contains reserved characters", ErrInvalidIdentity)
	}

	atIndex := strings.IndexByte(email, '@')
	if atIndex <= 0 || atIndex == len(email)-1 {
		return fmt.Errorf("%w: email must contain local and domain parts", ErrInvalidIdentity)
	}
	return nil
}

func validateIdentityTimestamp(timestamp string) error {
	if timestamp == "" {
		return fmt.Errorf("%w: timestamp is empty", ErrInvalidIdentity)
	}

	unix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: invalid timestamp %q", ErrInvalidIdentity, timestamp)
	}
	if unix < 0 {
		return fmt.Errorf("%w: timestamp must be non-negative", ErrInvalidIdentity)
	}
	return nil
}

func validateIdentityTimezone(timezone string) error {
	if len(timezone) != 5 {
		return fmt.Errorf("%w: timezone must be in ±HHMM format", ErrInvalidIdentity)
	}
	if timezone[0] != '+' && timezone[0] != '-' {
		return fmt.Errorf("%w: timezone must start with '+' or '-'", ErrInvalidIdentity)
	}

	hours, err := strconv.Atoi(timezone[1:3])
	if err != nil {
		return fmt.Errorf("%w: invalid timezone hour in %q", ErrInvalidIdentity, timezone)
	}
	if hours > 23 {
		return fmt.Errorf("%w: timezone hour out of range in %q", ErrInvalidIdentity, timezone)
	}

	minutes, err := strconv.Atoi(timezone[3:5])
	if err != nil {
		return fmt.Errorf("%w: invalid timezone minute in %q", ErrInvalidIdentity, timezone)
	}
	if minutes > 59 {
		return fmt.Errorf("%w: timezone minute out of range in %q", ErrInvalidIdentity, timezone)
	}
	return nil
}
