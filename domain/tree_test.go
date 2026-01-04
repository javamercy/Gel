package domain

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTreeEntry_Valid(t *testing.T) {
	hash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
	entry, err := NewTreeEntry(RegularFile, hash, "file.txt")
	assert.NoError(t, err)
	assert.Equal(t, RegularFile, entry.Mode)
	assert.Equal(t, hash, entry.Hash)
	assert.Equal(t, "file.txt", entry.Name)
}

func TestNewTreeEntry_Invalid(t *testing.T) {
	hash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	_, err := NewTreeEntry(RegularFile, "invalid-hash", "file.txt")
	assert.Error(t, err)

	_, err = NewTreeEntry(RegularFile, hash, "")
	assert.Error(t, err)
}

func TestTree_Deserialize_MultipleEntries(t *testing.T) {
	hash1 := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
	hash2 := "b1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	entry1, _ := NewTreeEntry(RegularFile, hash1, "file1.txt")
	entry2, _ := NewTreeEntry(Directory, hash2, "dir1")

	tree := NewTreeFromEntries([]TreeEntry{entry1, entry2})

	serialized := tree.Serialize()
	assert.Contains(t, string(serialized), "tree")

	newTree, err := NewTree(tree.Body())
	require.NoError(t, err)

	entries, err := newTree.Deserialize()
	require.NoError(t, err)
	assert.Len(t, entries, 2)

	assert.Equal(t, "file1.txt", entries[0].Name)
	assert.Equal(t, RegularFile, entries[0].Mode)
	assert.Equal(t, hash1, entries[0].Hash)

	assert.Equal(t, "dir1", entries[1].Name)
	assert.Equal(t, Directory, entries[1].Mode)
	assert.Equal(t, hash2, entries[1].Hash)
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
