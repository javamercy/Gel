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
		{0o100644, RegularFile},
		{0o100755, ExecutableFile},
		{0o040000, Directory},
		{0o120000, Symlink},
		{0o160000, Submodule},
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
		{"100644", RegularFile},
		{"100755", ExecutableFile},
		{"40000", Directory},
		{"120000", Symlink},
		{"160000", Submodule},
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
		{RegularFile, "100644"},
		{ExecutableFile, "100755"},
		{Directory, "40000"},
		{Symlink, "120000"},
		{Submodule, "160000"},
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
	assert.True(t, RegularFile.IsRegularFile())
	assert.False(t, Directory.IsRegularFile())

	assert.True(t, Directory.IsDirectory())
	assert.False(t, RegularFile.IsDirectory())

	assert.True(t, ExecutableFile.IsExecutableFile())
	assert.False(t, RegularFile.IsExecutableFile())

	assert.True(t, Symlink.IsSymlink())
	assert.False(t, RegularFile.IsSymlink())

	assert.True(t, Submodule.IsSubmodule())
	assert.False(t, RegularFile.IsSubmodule())

	assert.True(t, RegularFile.IsValid())
	assert.False(t, InvalidMode.IsValid())
}

func TestFileMode_ObjectType(t *testing.T) {
	tests := []struct {
		mode         FileMode
		expectedType ObjectType
		expectErr    bool
	}{
		{RegularFile, ObjectTypeBlob, false},
		{ExecutableFile, ObjectTypeBlob, false},
		{Symlink, ObjectTypeBlob, false},
		{Directory, ObjectTypeTree, false},
		{Submodule, ObjectTypeCommit, false},
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
