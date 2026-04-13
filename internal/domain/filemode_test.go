package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFileMode(t *testing.T) {
	tests := []struct {
		input    uint32
		expected FileMode
	}{
		{0o100644, FileModeRegular},
		{0o100755, FileModeExecutable},
		{0o040000, FileModeDirectory},
		{0, FileModeInvalid},
		{0o777, FileModeInvalid},
	}

	for _, tt := range tests {
		t.Run(
			"uint32", func(t *testing.T) {
				assert.Equal(t, tt.expected, ParseFileMode(tt.input))
			},
		)
	}
}

func TestParseFileModeFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected FileMode
	}{
		{"100644", FileModeRegular},
		{"100755", FileModeExecutable},
		{"40000", FileModeDirectory},
		{"invalid", FileModeInvalid},
		{"", FileModeInvalid},
	}

	for _, tt := range tests {
		t.Run(
			tt.input, func(t *testing.T) {
				assert.Equal(t, tt.expected, ParseFileModeFromString(tt.input))
			},
		)
	}
}

func TestFileMode_String(t *testing.T) {
	tests := []struct {
		mode     FileMode
		expected string
	}{
		{FileModeRegular, "100644"},
		{FileModeExecutable, "100755"},
		{FileModeDirectory, "40000"},
		{FileModeInvalid, ""},
	}

	for _, tt := range tests {
		t.Run(
			tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.mode.String())
			},
		)
	}
}

func TestFileMode_Helpers(t *testing.T) {
	assert.True(t, FileModeRegular.IsRegularFile())
	assert.False(t, FileModeDirectory.IsRegularFile())

	assert.True(t, FileModeDirectory.IsDirectory())
	assert.False(t, FileModeRegular.IsDirectory())

	assert.True(t, FileModeExecutable.IsExecutableFile())
	assert.False(t, FileModeRegular.IsExecutableFile())

	assert.True(t, FileModeRegular.IsValid())
	assert.False(t, FileModeInvalid.IsValid())
}

func TestFileMode_ObjectType(t *testing.T) {
	tests := []struct {
		mode         FileMode
		expectedType ObjectType
		expectErr    bool
	}{
		{FileModeRegular, ObjectTypeBlob, false},
		{FileModeExecutable, ObjectTypeBlob, false},
		{FileModeDirectory, ObjectTypeTree, false},
		{FileModeInvalid, "", true},
	}

	for _, tt := range tests {
		t.Run(
			tt.mode.String(), func(t *testing.T) {
				objType, err := tt.mode.ObjectType()
				if tt.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedType, objType)
				}
			},
		)
	}
}

func TestParseFileModeFromOsMode(t *testing.T) {
	assert.Equal(t, FileModeDirectory, ParseFileModeFromOsMode(0o040755))
	assert.Equal(t, FileModeExecutable, ParseFileModeFromOsMode(0o100755))
	assert.Equal(t, FileModeRegular, ParseFileModeFromOsMode(0o100644))
}

func TestFileMode_Utils(t *testing.T) {
	assert.Equal(t, uint32(0o100644), FileModeRegular.Uint32())
	assert.True(t, FileModeRegular.Equals(FileModeRegular))
	assert.False(t, FileModeRegular.Equals(FileModeExecutable))
}
