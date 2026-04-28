package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmptyIndex(t *testing.T) {
	index := NewEmptyIndex()

	assert.Equal(t, IndexSignature, string(index.Header.Signature[:]))
	assert.Equal(t, uint32(IndexVersion), index.Header.Version)
	assert.Equal(t, uint32(0), index.Header.NumEntries)
	assert.Empty(t, index.Entries)
}

func TestAddEntry_SingleEntry(t *testing.T) {
	index := NewEmptyIndex()
	entry := createTestEntry("a.txt", "hash1")

	index.AddEntry(entry)

	assert.Equal(t, 1, len(index.Entries))
	assert.Equal(t, uint32(1), index.Header.NumEntries)
	assert.Equal(t, entry, index.Entries[0])
}

func TestAddEntry_MultipleEntries(t *testing.T) {
	index := NewEmptyIndex()
	entry1 := createTestEntry("a.txt", "hash1")
	entry2 := createTestEntry("b.txt", "hash2")
	entry3 := createTestEntry("c.txt", "hash3")

	index.AddEntry(entry1)
	index.AddEntry(entry2)
	index.AddEntry(entry3)

	assert.Equal(t, 3, len(index.Entries))
	assert.Equal(t, uint32(3), index.Header.NumEntries)
}

func TestAddOrUpdateEntry_NewPath(t *testing.T) {
	index := NewEmptyIndex()
	entry := createTestEntry("new.txt", "hash1")

	index.SetEntry(entry)

	assert.Equal(t, 1, len(index.Entries))
	assert.Equal(t, NormalizedPath("new.txt"), index.Entries[0].Path)
}

func TestAddOrUpdateEntry_ExistingPath(t *testing.T) {
	index := NewEmptyIndex()
	entry1 := createTestEntry("a.txt", "hash1")
	index.AddEntry(entry1)

	entry2 := createTestEntry("a.txt", "hash2")
	index.SetEntry(entry2)

	assert.Equal(t, 1, len(index.Entries))
	assert.Equal(t, entry2.Hash, index.Entries[0].Hash)
}

func TestRemoveEntry_ExistingPath(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))
	index.AddEntry(createTestEntry("b.txt", "hash2"))
	index.AddEntry(createTestEntry("c.txt", "hash3"))

	index.RemoveEntry("b.txt")

	assert.Equal(t, 2, len(index.Entries))
	assert.Equal(t, uint32(2), index.Header.NumEntries)
	assert.False(t, index.HasEntry("b.txt"))
}

func TestRemoveEntry_NonExistentPath(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))
	index.AddEntry(createTestEntry("b.txt", "hash2"))

	index.RemoveEntry("nonexistent.txt")

	assert.Equal(t, 2, len(index.Entries))
	assert.Equal(t, uint32(2), index.Header.NumEntries)
}

func TestFindEntry_Exists(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))
	index.AddEntry(createTestEntry("b.txt", "hash2"))
	index.AddEntry(createTestEntry("c.txt", "hash3"))

	entry, _ := index.FindEntry("b.txt")

	require.NotNil(t, entry)
	assert.Equal(t, NormalizedPath("b.txt"), entry.Path)
}

func TestUpdateEntry_NotFound(t *testing.T) {
	index := NewEmptyIndex()
	updated := index.UpdateEntry(createTestEntry("missing.txt", "hash1"))
	assert.False(t, updated)
}

func TestFindEntriesByPathPrefix(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("src/a.go", "hash1"))
	index.AddEntry(createTestEntry("src/b.go", "hash2"))
	index.AddEntry(createTestEntry("docs/readme.md", "hash3"))

	entries := index.FindEntriesByPathPrefix("src/")
	require.Len(t, entries, 2)
	assert.Equal(t, NormalizedPath("src/a.go"), entries[0].Path)
	assert.Equal(t, NormalizedPath("src/b.go"), entries[1].Path)
}

func TestFindEntriesByPathPattern(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("src/a.go", "hash1"))
	index.AddEntry(createTestEntry("src/b.txt", "hash2"))
	index.AddEntry(createTestEntry("src/c.go", "hash3"))

	entries := index.FindEntriesByPathPattern("src/*.go")
	require.Len(t, entries, 2)
	assert.Equal(t, NormalizedPath("src/a.go"), entries[0].Path)
	assert.Equal(t, NormalizedPath("src/c.go"), entries[1].Path)
}

func TestComputeIndexFlagsAndGetStage(t *testing.T) {
	flags := ComputeIndexFlags("path/to/file.txt", 2)
	entry := createTestEntry("path/to/file.txt", "hash1")
	entry.Flags = flags

	assert.Equal(t, uint16(2), entry.GetStage())
	assert.Equal(t, uint16(len("path/to/file.txt")), flags&MaxPathLength)
}

func TestFindEntry_NotExists(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))

	entry, _ := index.FindEntry("nonexistent.txt")

	assert.Nil(t, entry)
}

func TestHasEntry_True(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))
	index.AddEntry(createTestEntry("b.txt", "hash2"))

	assert.True(t, index.HasEntry("a.txt"))
	assert.True(t, index.HasEntry("b.txt"))
}

func TestHasEntry_False(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))

	assert.False(t, index.HasEntry("nonexistent.txt"))
	assert.False(t, index.HasEntry(""))
}

func TestSerialize_EmptyIndex(t *testing.T) {
	index := NewEmptyIndex()

	data, err := index.Serialize()

	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Equal(t, 44, len(data))
}

func TestSerialize_WithEntries(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))
	index.AddEntry(createTestEntry("b.txt", "hash2"))

	data, err := index.Serialize()

	require.NoError(t, err)
	assert.Greater(t, len(data), 44)
}

func TestDeserialize_Valid(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("test.txt", "hash1"))
	index.AddEntry(createTestEntry("another.txt", "hash2"))

	data, err := index.Serialize()
	require.NoError(t, err)

	deserializedIndex, err := DeserializeIndex(data)

	require.NoError(t, err)
	assert.Equal(t, index.Header.NumEntries, deserializedIndex.Header.NumEntries)
	assert.Equal(t, len(index.Entries), len(deserializedIndex.Entries))
}

func TestDeserialize_InvalidSignature(t *testing.T) {
	invalidData := []byte("XXXX_INVALID_DATA")

	_, err := DeserializeIndex(invalidData)

	assert.ErrorIs(t, err, ErrInvalidIndexSignature)
}

func TestDeserialize_ChecksumMismatch(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))

	data, err := index.Serialize()
	require.NoError(t, err)

	data[len(data)-1] ^= 0xFF

	_, err = DeserializeIndex(data)

	assert.ErrorIs(t, err, ErrChecksumMismatch)
}

func TestSerializeDeserializeIndex_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
	}{
		{"empty", []string{}},
		{"single", []string{"test.txt"}},
		{"multiple", []string{"a.txt", "b.txt", "c.txt"}},
		{"nested", []string{"dir/file.txt", "dir/sub/deep.txt"}},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				index := NewEmptyIndex()
				for i, path := range tt.paths {
					index.AddEntry(createTestEntry(path, fmt.Sprintf("hash%d", i)))
				}

				data, err := index.Serialize()
				require.NoError(t, err)

				deserializedIndex, err := DeserializeIndex(data)
				require.NoError(t, err)

				assert.Equal(t, index.Header.NumEntries, deserializedIndex.Header.NumEntries)
				assert.Equal(t, len(index.Entries), len(deserializedIndex.Entries))

				for i := range index.Entries {
					assert.Equal(t, index.Entries[i].Path, deserializedIndex.Entries[i].Path)
					assert.Equal(t, index.Entries[i].Hash, deserializedIndex.Entries[i].Hash)
				}
			},
		)
	}
}

func TestDeserialize_EmptyData(t *testing.T) {
	var data []byte
	index, err := DeserializeIndex(data)
	require.NoError(t, err)
	assert.Equal(t, uint32(0), index.Header.NumEntries)
}

func TestDeserialize_TruncatedData(t *testing.T) {
	truncatedData := []byte{0x44, 0x49, 0x52, 0x43}
	_, err := DeserializeIndex(truncatedData)
	assert.Error(t, err)
}

func TestDeserialize_IncorrectChecksumSize(t *testing.T) {
	index := NewEmptyIndex()
	index.AddEntry(createTestEntry("a.txt", "hash1"))

	data, err := index.Serialize()
	require.NoError(t, err)

	broken := data[:len(data)-1]
	_, err = DeserializeIndex(broken)
	assert.ErrorIs(t, err, ErrIncorrectChecksumSize)
}

func TestIndexEntry_MatchesStat(t *testing.T) {
	entry := createTestEntry("a.txt", "hash1")
	entry.Mode = ParseFileModeFromOsMode(0o100644).Uint32()
	stat := &FileStat{
		Device:       entry.Device,
		Inode:        entry.Inode,
		Mode:         0o100644,
		Size:         entry.Size,
		ChangedTime:  entry.ChangedTime,
		ModifiedTime: entry.ModifiedTime,
	}

	assert.True(t, entry.MatchesStat(stat))

	stat.Size++
	assert.False(t, entry.MatchesStat(stat))
}

func createTestEntry(pathStr, hashSeed string) *IndexEntry {
	path, err := NewNormalizedPathUnchecked(pathStr)
	if err != nil {
		panic(err)
	}
	fullHash, err := NewHash(ComputeSHA256([]byte(hashSeed)))
	if err != nil {
		panic(err)
	}
	entry := NewIndexEntry(
		path,
		fullHash,
		100,
		uint32(FileModeRegular),
		0, 0, 0, 0,
		uint16(len(pathStr)),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	return entry
}
