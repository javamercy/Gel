package domain

import (
	"Gel/core/constant"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var hashList = []string{
	"625786dc2c6ad841446f9394fd439ce53bf6442354d0b0d095b97f6b12499419",
	"4c64cbb574119e68e31489ac8f38350b43f3ac75361a472563dd468f70b77652",
	"658574d668dec4d0bb7145fed0ab97a3c348f0db21389e65416c03557ea48c3d",
}

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
	entry := createTestEntry("a.txt", hashList[0])

	index.AddEntry(entry)
	assert.Equal(t, 1, len(index.Entries), "entries slice should have length 1")
	assert.Equal(t, uint32(len(index.Entries)), index.Header.NumEntries, "should have 1 entry")

	addedEntry := index.Entries[0]
	assert.Equal(t, addedEntry, entry)
}

func TestAddEntry_MultipleEntries(t *testing.T) {
	index := NewEmptyIndex()
	firstEntry := createTestEntry("a.txt", hashList[0])
	secondEntry := createTestEntry("b.txt", hashList[1])

	index.AddEntry(firstEntry)
	index.AddEntry(secondEntry)

	assert.Equal(t, 2, len(index.Entries), "entries len should be 2")
	assert.Equal(t, uint32(len(index.Entries)), index.Header.NumEntries, "entries slice length should match numEntries")

	firstAddedEntry := index.Entries[0]
	secondAddedEntry := index.Entries[1]

	assert.Equal(t, firstEntry, firstAddedEntry)
	assert.Equal(t, secondEntry, secondAddedEntry)
}

func TestAddOrUpdateEntry_ExistingPath(t *testing.T) {
	index := NewEmptyIndex()
	entry := createTestEntry("a.txt", hashList[0])

	index.AddEntry(entry)
	assert.Equal(t, entry, index.Entries[0])

	entry.Hash = hashList[1]

	index.AddOrUpdateEntry(entry)
	assert.Equal(t, entry, index.Entries[0])
}

func TestRemoveEntry_ExistingPath(t *testing.T) {
	index := NewEmptyIndex()
	entry := createTestEntry("a.txt", hashList[0])

	index.AddEntry(entry)
	assert.Equal(t, 1, len(index.Entries))
	assert.Equal(t, index.Header.NumEntries, uint32(len(index.Entries)))

	index.RemoveEntry(entry.Path)
	assert.Equal(t, 0, len(index.Entries))
	assert.Equal(t, index.Header.NumEntries, uint32(len(index.Entries)))
}

func TestRemoveEntry_NonExistentPath(t *testing.T) {
	index := NewEmptyIndex()
	entry := createTestEntry("a.txt", hashList[0])

	index.AddEntry(entry)
	assert.Equal(t, 1, len(index.Entries))
	assert.Equal(t, index.Header.NumEntries, uint32(len(index.Entries)))

	index.RemoveEntry("nonexistent path")
	assert.Equal(t, 1, len(index.Entries))
	assert.Equal(t, index.Header.NumEntries, uint32(len(index.Entries)))
}

// func TestSerialize_EmptyIndex(t *testing.T)
// func TestDeserialize_InvalidSignature(t *testing.T)
// func TestDeserialize_ChecksumMismatch(t *testing.T)

func createTestEntry(path, hash string) *IndexEntry {
	return NewIndexEntry(
		path,
		hash,
		100,
		uint32(RegularFile),
		0, 0, 0, 0,
		uint16(len(path)),
		time.Now(),
		time.Now(),
	)
}
