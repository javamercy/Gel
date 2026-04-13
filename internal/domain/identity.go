package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidIdentityFormat is returned when identity input violates required shape or validation rules.
var ErrInvalidIdentityFormat = errors.New("invalid identity format")

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

	if name == "" || email == "" {
		return Identity{}, ErrInvalidIdentityFormat
	}
	if !isValidEmail(email) {
		return Identity{}, ErrInvalidIdentityFormat
	}
	if _, err := strconv.ParseInt(timestamp, 10, 64); err != nil {
		return Identity{}, ErrInvalidIdentityFormat
	}
	if !isValidTimezone(timezone) {
		return Identity{}, ErrInvalidIdentityFormat
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
func (identity Identity) Serialize() []byte {
	return []byte(fmt.Sprintf(
		"%s <%s> %s %s",
		identity.Name,
		identity.Email,
		identity.Timestamp,
		identity.Timezone,
	))
}

// isValidEmail reports whether email satisfies the minimal format expected by domain validation.
func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}

// isValidTimezone reports whether timezone is in valid ±HHMM format and range.
func isValidTimezone(timezone string) bool {
	if len(timezone) != 5 {
		return false
	}
	if timezone[0] != '+' && timezone[0] != '-' {
		return false
	}

	hh, err := strconv.Atoi(timezone[1:3])
	if err != nil || hh < 0 || hh > 23 {
		return false
	}

	mm, err := strconv.Atoi(timezone[3:5])
	if err != nil || mm < 0 || mm > 59 {
		return false
	}
	return true
}
