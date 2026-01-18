package util

import (
	"strconv"
	"time"
)

func FormatCommitTimestamp(time time.Time) string {
	return strconv.FormatInt(time.Unix(), 10)
}

func FormatCommitTimezone(time time.Time) string {
	return time.Format("-0700")
}

func FormatCommitDate(timestamp, timezone string) (string, error) {
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
