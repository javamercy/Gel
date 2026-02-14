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
		{0o100644, RegularFileMode},
		{0o100755, ExecutableFileMode},
		{0o040000, DirectoryMode},
		{0, InvalidMode},
		{0o777, InvalidMode},
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
		{"100644", RegularFileMode},
		{"100755", ExecutableFileMode},
		{"40000", DirectoryMode},
		{"invalid", InvalidMode},
		{"", InvalidMode},
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
		{RegularFileMode, "100644"},
		{ExecutableFileMode, "100755"},
		{DirectoryMode, "40000"},
		{InvalidMode, ""},
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
	assert.True(t, RegularFileMode.IsRegularFile())
	assert.False(t, DirectoryMode.IsRegularFile())

	assert.True(t, DirectoryMode.IsDirectory())
	assert.False(t, RegularFileMode.IsDirectory())

	assert.True(t, ExecutableFileMode.IsExecutableFile())
	assert.False(t, RegularFileMode.IsExecutableFile())

	assert.True(t, RegularFileMode.IsValid())
	assert.False(t, InvalidMode.IsValid())
}

func TestFileMode_ObjectType(t *testing.T) {
	tests := []struct {
		mode         FileMode
		expectedType ObjectType
		expectErr    bool
	}{
		{RegularFileMode, ObjectTypeBlob, false},
		{ExecutableFileMode, ObjectTypeBlob, false},
		{DirectoryMode, ObjectTypeTree, false},
		{InvalidMode, "", true},
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
