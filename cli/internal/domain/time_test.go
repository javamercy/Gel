package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatCommitTimestampAndTimezone(t *testing.T) {
	loc := time.FixedZone("+0300", 3*60*60)
	ts := time.Unix(1710000000, 0).In(loc)

	assert.Equal(t, "1710000000", FormatCommitTimestamp(ts))
	assert.Equal(t, "+0300", FormatCommitTimezone(ts))
}

func TestFormatCommitDate(t *testing.T) {
	out, err := FormatCommitDate("0", "+0300")
	require.NoError(t, err)
	assert.Equal(t, "1970-01-01 03:00:00 +0300", out)
}

func TestFormatCommitDate_InvalidInputs(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
		timezone  string
	}{
		{name: "empty timestamp", timestamp: "", timezone: "+0000"},
		{name: "bad timezone sign", timestamp: "1", timezone: "0000"},
		{name: "bad timezone len", timestamp: "1", timezone: "+000"},
		{name: "bad timestamp", timestamp: "abc", timezone: "+0000"},
		{name: "bad timezone hour", timestamp: "1", timezone: "+xx00"},
		{name: "bad timezone minute", timestamp: "1", timezone: "+00xx"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := FormatCommitDate(tt.timestamp, tt.timezone)
				assert.Error(t, err)
			},
		)
	}
}
