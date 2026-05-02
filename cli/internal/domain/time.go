package domain

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrInvalidCommitTimestamp = errors.New("invalid commit timestamp")
	ErrInvalidCommitTimezone  = errors.New("invalid commit timezone")
)

// FormatCommitTimestamp formats a time as Unix timestamp string for commit objects.
func FormatCommitTimestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

// FormatCommitTimezone formats a time's timezone in Git's timezone format (e.g., "+0530", "-0700").
func FormatCommitTimezone(t time.Time) string {
	return t.Format("-0700")
}

// FormatCommitDate parses a commit timestamp and timezone and returns a
// formatted date string in "2006-01-02 15:04:05 -0700" layout.
// It validates that timestamp is non-empty and timezone is in ±HHMM format.
func FormatCommitDate(timestamp, timezone string) (string, error) {
	unix, err := parseCommitTimestamp(timestamp)
	if err != nil {
		return "", err
	}

	offset, err := parseCommitTimezoneOffset(timezone)
	if err != nil {
		return "", err
	}

	loc := time.FixedZone(timezone, offset)
	return time.Unix(unix, 0).In(loc).Format("2006-01-02 15:04:05 -0700"), nil
}

func parseCommitTimestamp(timestamp string) (int64, error) {
	if timestamp == "" {
		return 0, fmt.Errorf("%w: empty", ErrInvalidCommitTimestamp)
	}

	unix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidCommitTimestamp, timestamp)
	}
	return unix, nil
}

func parseCommitTimezoneOffset(timezone string) (int, error) {
	if len(timezone) != 5 || (timezone[0] != '+' && timezone[0] != '-') {
		return 0, fmt.Errorf("%w: %q must be ±HHMM", ErrInvalidCommitTimezone, timezone)
	}
	hours, err := strconv.Atoi(timezone[1:3])
	if err != nil {
		return 0, fmt.Errorf("%w: invalid hour in %q", ErrInvalidCommitTimezone, timezone)
	}
	minutes, err := strconv.Atoi(timezone[3:5])
	if err != nil {
		return 0, fmt.Errorf("%w: invalid minute in %q", ErrInvalidCommitTimezone, timezone)
	}
	if hours > 23 || minutes > 59 {
		return 0, fmt.Errorf("%w: %q out of range", ErrInvalidCommitTimezone, timezone)
	}
	offset := (hours*60 + minutes) * 60
	if timezone[0] == '-' {
		offset = -offset
	}
	return offset, nil
}
