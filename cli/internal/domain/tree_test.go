package domain

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTreeEntry_Valid(t *testing.T) {
	h, err := NewHash("a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	require.NoError(t, err)
	entry := NewTreeEntry(FileModeRegular, h, "file.txt")
	assert.Equal(t, FileModeRegular, entry.Mode)
	assert.Equal(t, h, entry.Hash)
	assert.Equal(t, "file.txt", entry.Name)
}

func TestTree_Deserialize_MultipleEntries(t *testing.T) {
	h1, err := NewHash("a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	require.NoError(t, err)
	h2, err := NewHash("b1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	require.NoError(t, err)

	entry1 := NewTreeEntry(FileModeRegular, h1, "file1.txt")
	entry2 := NewTreeEntry(FileModeDirectory, h2, "dir1")

	tree, err := NewTreeFromEntries([]TreeEntry{entry1, entry2})
	require.NoError(t, err)

	serialized := tree.Serialize()
	assert.Contains(t, string(serialized), "tree")

	newTree, err := NewTree(tree.Body())
	require.NoError(t, err)

	entries, err := newTree.Deserialize()
	require.NoError(t, err)
	assert.Len(t, entries, 2)

	assert.Equal(t, "file1.txt", entries[0].Name)
	assert.Equal(t, FileModeRegular, entries[0].Mode)
	assert.Equal(t, h1, entries[0].Hash)

	assert.Equal(t, "dir1", entries[1].Name)
	assert.Equal(t, FileModeDirectory, entries[1].Mode)
	assert.Equal(t, h2, entries[1].Hash)
}

func TestTree_Deserialize_InvalidFormat(t *testing.T) {
	hashBytes, _ := hex.DecodeString("a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")

	body := []byte("100644 file.txt")
	body = append(body, hashBytes...)

	_, err := NewTree(body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing null byte")

	body2 := []byte("invalid_mode file.txt\x00")
	body2 = append(body2, hashBytes...)

	_, err = NewTree(body2)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidFileMode, err)
}

func TestNewTreeFromEntries_InvalidNameOrMode(t *testing.T) {
	h, err := NewHash("a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	require.NoError(t, err)

	_, err = NewTreeFromEntries([]TreeEntry{NewTreeEntry(FileModeRegular, h, "")})
	assert.ErrorIs(t, err, ErrInvalidTreeEntryName)

	_, err = NewTreeFromEntries([]TreeEntry{NewTreeEntry(FileModeInvalid, h, "file.txt")})
	assert.ErrorIs(t, err, ErrInvalidFileMode)
}

func TestTree_BodyIsDefensiveCopy(t *testing.T) {
	h, err := NewHash("a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	require.NoError(t, err)
	tree, err := NewTreeFromEntries([]TreeEntry{NewTreeEntry(FileModeRegular, h, "file.txt")})
	require.NoError(t, err)

	body := tree.Body()
	body[0] = '9'
	assert.NotEqual(t, body[0], tree.Body()[0])
}
