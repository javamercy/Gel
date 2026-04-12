package domain

import (
	"errors"
	"strconv"
	"time"
)

// FormatCommitTimestamp formats a time as Unix timestamp string for commit objects.
func FormatCommitTimestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

// FormatCommitTimezone formats a time's timezone in Git's timezone format (e.g., "+0530", "-0700").
func FormatCommitTimezone(t time.Time) string {
	return t.Format("-0700")
}

// FormatCommitDate parses a commit timestamp and timezone and returns a formatted date string.
// It validates that timestamp is non-empty and timezone is in ±HHMM format.
func FormatCommitDate(timestamp, timezone string) (string, error) {
	if len(timestamp) == 0 {
		return "", errors.New("invalid timestamp: empty")
	}
	if len(timezone) != 5 || (timezone[0] != '+' && timezone[0] != '-') {
		return "", errors.New("invalid timezone: must be ±HHMM format")
	}

	unix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "", err
	}

	hours, err := strconv.Atoi(timezone[1:3])
	if err != nil {
		return "", err
	}

	mins, err := strconv.Atoi(timezone[3:5])
	if err != nil {
		return "", err
	}

	offset := (hours*60 + mins) * 60
	if timezone[0] == '-' {
		offset = -offset
	}

	loc := time.FixedZone(timezone, offset)
	t := time.Unix(unix, 0).In(loc)
	return t.Format("2006-01-02 15:04:05 -0700"), nil
}
