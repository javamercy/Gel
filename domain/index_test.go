package domain

import (
	"Gel/core/constant"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptyIndex(t *testing.T) {
	index := NewEmptyIndex()

	assert.Equal(t, constant.GelIndexSignature, string(index.Header.Signature[:]), "signature should be DIRC")
	assert.Equal(t, uint32(constant.GelIndexVersion), index.Header.Version, "version should be 1")
	assert.Equal(t, uint32(0), index.Header.NumEntries, "empty index should have 0 entries")

	assert.NotNil(t, index.Entries, "entries slice should not be nil")
	assert.Empty(t, index.Entries, "entries slice should be empty")

	assert.Equal(t, "", index.Checksum, "checksum should be empty string")
}
func TestAddEntry_SingleEntry(t *testing.T) {
	index := NewEmptyIndex()
	entry := NewIndexEntry(
		"test.txt",
		"abc123def456",
		100,
		uint32(RegularFile),
		0, 0, 0, 0,
		8,
		time.Now(),
		time.Now(),
	)

	index.AddEntry(entry)
	assert.Equal(t, uint32(1), index.Header.NumEntries, "should have 1 entry")
	assert.Len(t, index.Entries, 1, "entries slice should have length 1")

	addedEntry := index.Entries[0]
	assert.ElementsMatch(t, addedEntry, entry)
}
func TestAddEntry_MultipleEntries(t *testing.T) {
	index := NewEmptyIndex()
	firstEntry := NewIndexEntry(
		"test1.txt",
		"123456789",
		100,
		uint32(RegularFile),
		0, 0, 0, 0,
		8,
		time.Now(),
		time.Now(),
	)
	secondEntry := NewIndexEntry(
		"test2.txt",
		"123654789",
		100,
		uint32(RegularFile),
		0, 0, 0, 0,
		8,
		time.Now(),
		time.Now(),
	)

	index.AddEntry(firstEntry)
	index.AddEntry(secondEntry)

	assert.Equal(t, 2, len(index.Entries), "entries len should be 2")
	assert.Equal(t, len(index.Entries), index.Header.NumEntries, "entries slice length should match numEntries")

	firstAddedEntry := index.Entries[0]
	secondAddedEntry := index.Entries[1]

	assert.Equal(t, firstEntry, firstAddedEntry)
	assert.Equal(t, secondEntry, secondAddedEntry)
}

// func TestAddOrUpdateEntry_NewPath(t *testing.T)
// func TestAddOrUpdateEntry_ExistingPath(t *testing.T)
// func TestRemoveEntry_ExistingPath(t *testing.T)
// func TestRemoveEntry_NonExistentPath(t *testing.T)
// func TestSerialize_EmptyIndex(t *testing.T)
// func TestDeserialize_InvalidSignature(t *testing.T)
// func TestDeserialize_ChecksumMismatch(t *testing.T)
