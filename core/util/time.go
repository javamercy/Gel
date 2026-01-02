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
