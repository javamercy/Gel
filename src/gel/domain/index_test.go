package domain

import (
	"testing"
	"time"
)

// TestNewIndex tests the NewIndex constructor
func TestNewIndex(t *testing.T) {
	tests := []struct {
		name     string
		header   IndexHeader
		entries  []IndexEntry
		checksum string
	}{
		{
			name: "create index with valid header and entries",
			header: IndexHeader{
				Signature:  [4]byte{'G', 'E', 'L', 'I'},
				Version:    1,
				NumEntries: 2,
			},
			entries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
				{Path: "file2.txt", Hash: "def456"},
			},
			checksum: "checksum123",
		},
		{
			name: "create index with empty entries",
			header: IndexHeader{
				Signature:  [4]byte{'G', 'E', 'L', 'I'},
				Version:    1,
				NumEntries: 0,
			},
			entries:  []IndexEntry{},
			checksum: "",
		},
		{
			name: "create index with nil checksum",
			header: IndexHeader{
				Signature:  [4]byte{'T', 'E', 'S', 'T'},
				Version:    2,
				NumEntries: 1,
			},
			entries: []IndexEntry{
				{Path: "test.txt", Hash: "hash1"},
			},
			checksum: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := NewIndex(tt.header, tt.entries, tt.checksum)

			if index == nil {
				t.Fatal("NewIndex returned nil")
			}

			if index.Header != tt.header {
				t.Errorf("Header = %v, want %v", index.Header, tt.header)
			}

			if len(index.Entries) != len(tt.entries) {
				t.Errorf("Entries length = %d, want %d", len(index.Entries), len(tt.entries))
			}

			if index.Checksum != tt.checksum {
				t.Errorf("Checksum = %s, want %s", index.Checksum, tt.checksum)
			}
		})
	}
}

// TestNewEmptyIndex tests the NewEmptyIndex constructor
func TestNewEmptyIndex(t *testing.T) {
	index := NewEmptyIndex()

	if index == nil {
		t.Fatal("NewEmptyIndex returned nil")
	}

	expectedSignature := [4]byte{'G', 'E', 'L', 'I'}
	if index.Header.Signature != expectedSignature {
		t.Errorf("Signature = %v, want %v", index.Header.Signature, expectedSignature)
	}

	if index.Header.Version != 1 {
		t.Errorf("Version = %d, want 1", index.Header.Version)
	}

	if index.Header.NumEntries != 0 {
		t.Errorf("NumEntries = %d, want 0", index.Header.NumEntries)
	}

	if len(index.Entries) != 0 {
		t.Errorf("Entries length = %d, want 0", len(index.Entries))
	}

	if index.Checksum != "" {
		t.Errorf("Checksum = %s, want empty string", index.Checksum)
	}
}

// TestAddEntry tests the AddEntry method
func TestAddEntry(t *testing.T) {
	tests := []struct {
		name           string
		initialEntries []IndexEntry
		entriesToAdd   []IndexEntry
		expectedCount  int
		expectedPaths  []string
	}{
		{
			name:           "add entry to empty index",
			initialEntries: []IndexEntry{},
			entriesToAdd: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			expectedCount: 1,
			expectedPaths: []string{"file1.txt"},
		},
		{
			name: "add entry to non-empty index",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			entriesToAdd: []IndexEntry{
				{Path: "file2.txt", Hash: "def456"},
			},
			expectedCount: 2,
			expectedPaths: []string{"file1.txt", "file2.txt"},
		},
		{
			name:           "add multiple entries",
			initialEntries: []IndexEntry{},
			entriesToAdd: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
				{Path: "file2.txt", Hash: "def456"},
				{Path: "file3.txt", Hash: "ghi789"},
			},
			expectedCount: 3,
			expectedPaths: []string{"file1.txt", "file2.txt", "file3.txt"},
		},
		{
			name:           "add entry with same path (duplicate - allowed in AddEntry)",
			initialEntries: []IndexEntry{},
			entriesToAdd: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
				{Path: "file1.txt", Hash: "xyz999"},
			},
			expectedCount: 2,
			expectedPaths: []string{"file1.txt", "file1.txt"},
		},
		{
			name:           "add entry with complete metadata",
			initialEntries: []IndexEntry{},
			entriesToAdd: []IndexEntry{
				{
					Path:        "test.txt",
					Hash:        "hash123",
					Size:        1024,
					Mode:        0644,
					Device:      1,
					Inode:       12345,
					UserId:      1000,
					GroupId:     1000,
					Flags:       0,
					CreatedTime: time.Now(),
					UpdatedTime: time.Now(),
				},
			},
			expectedCount: 1,
			expectedPaths: []string{"test.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := IndexHeader{
				Signature:  [4]byte{'G', 'E', 'L', 'I'},
				Version:    1,
				NumEntries: uint32(len(tt.initialEntries)),
			}
			index := NewIndex(header, tt.initialEntries, "")

			for _, entry := range tt.entriesToAdd {
				index.AddEntry(entry)
			}

			if index.Header.NumEntries != uint32(tt.expectedCount) {
				t.Errorf("NumEntries = %d, want %d", index.Header.NumEntries, tt.expectedCount)
			}

			if len(index.Entries) != tt.expectedCount {
				t.Errorf("Entries length = %d, want %d", len(index.Entries), tt.expectedCount)
			}

			for i, expectedPath := range tt.expectedPaths {
				if index.Entries[i].Path != expectedPath {
					t.Errorf("Entry[%d].Path = %s, want %s", i, index.Entries[i].Path, expectedPath)
				}
			}
		})
	}
}

// TestAddOrUpdateEntry tests the AddOrUpdateEntry method
func TestAddOrUpdateEntry(t *testing.T) {
	tests := []struct {
		name            string
		initialEntries  []IndexEntry
		entriesToAdd    []IndexEntry
		expectedCount   int
		expectedEntries []IndexEntry
	}{
		{
			name:           "add new entry to empty index",
			initialEntries: []IndexEntry{},
			entriesToAdd: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			expectedCount: 1,
			expectedEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
		},
		{
			name: "update existing entry",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123", Size: 100},
			},
			entriesToAdd: []IndexEntry{
				{Path: "file1.txt", Hash: "xyz999", Size: 200},
			},
			expectedCount: 1,
			expectedEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "xyz999", Size: 200},
			},
		},
		{
			name: "add new entry to non-empty index",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			entriesToAdd: []IndexEntry{
				{Path: "file2.txt", Hash: "def456"},
			},
			expectedCount: 2,
			expectedEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
				{Path: "file2.txt", Hash: "def456"},
			},
		},
		{
			name: "update entry then add new entry",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			entriesToAdd: []IndexEntry{
				{Path: "file1.txt", Hash: "updated123"},
				{Path: "file2.txt", Hash: "new456"},
			},
			expectedCount: 2,
			expectedEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "updated123"},
				{Path: "file2.txt", Hash: "new456"},
			},
		},
		{
			name: "update multiple times (no duplicates)",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "v1"},
			},
			entriesToAdd: []IndexEntry{
				{Path: "file1.txt", Hash: "v2"},
				{Path: "file1.txt", Hash: "v3"},
				{Path: "file1.txt", Hash: "v4"},
			},
			expectedCount: 1,
			expectedEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "v4"},
			},
		},
		{
			name: "update entry in middle of list",
			initialEntries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
				{Path: "c.txt", Hash: "hash_c"},
			},
			entriesToAdd: []IndexEntry{
				{Path: "b.txt", Hash: "hash_b_updated"},
			},
			expectedCount: 3,
			expectedEntries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b_updated"},
				{Path: "c.txt", Hash: "hash_c"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := IndexHeader{
				Signature:  [4]byte{'G', 'E', 'L', 'I'},
				Version:    1,
				NumEntries: uint32(len(tt.initialEntries)),
			}
			index := NewIndex(header, tt.initialEntries, "")

			for _, entry := range tt.entriesToAdd {
				index.AddOrUpdateEntry(entry)
			}

			if index.Header.NumEntries != uint32(tt.expectedCount) {
				t.Errorf("NumEntries = %d, want %d", index.Header.NumEntries, tt.expectedCount)
			}

			if len(index.Entries) != tt.expectedCount {
				t.Errorf("Entries length = %d, want %d", len(index.Entries), tt.expectedCount)
			}

			for i, expectedEntry := range tt.expectedEntries {
				if index.Entries[i].Path != expectedEntry.Path {
					t.Errorf("Entry[%d].Path = %s, want %s", i, index.Entries[i].Path, expectedEntry.Path)
				}
				if index.Entries[i].Hash != expectedEntry.Hash {
					t.Errorf("Entry[%d].Hash = %s, want %s", i, index.Entries[i].Hash, expectedEntry.Hash)
				}
			}
		})
	}
}

// TestRemoveEntry tests the RemoveEntry method
func TestRemoveEntry(t *testing.T) {
	tests := []struct {
		name           string
		initialEntries []IndexEntry
		pathsToRemove  []string
		expectedCount  int
		expectedPaths  []string
	}{
		{
			name: "remove existing entry",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
				{Path: "file2.txt", Hash: "def456"},
			},
			pathsToRemove: []string{"file1.txt"},
			expectedCount: 1,
			expectedPaths: []string{"file2.txt"},
		},
		{
			name: "remove non-existent entry",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			pathsToRemove: []string{"nonexistent.txt"},
			expectedCount: 1,
			expectedPaths: []string{"file1.txt"},
		},
		{
			name:           "remove from empty index",
			initialEntries: []IndexEntry{},
			pathsToRemove:  []string{"file1.txt"},
			expectedCount:  0,
			expectedPaths:  []string{},
		},
		{
			name: "remove all entries",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
				{Path: "file2.txt", Hash: "def456"},
			},
			pathsToRemove: []string{"file1.txt", "file2.txt"},
			expectedCount: 0,
			expectedPaths: []string{},
		},
		{
			name: "remove entry from middle",
			initialEntries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
				{Path: "c.txt", Hash: "hash_c"},
			},
			pathsToRemove: []string{"b.txt"},
			expectedCount: 2,
			expectedPaths: []string{"a.txt", "c.txt"},
		},
		{
			name: "remove first entry",
			initialEntries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
			},
			pathsToRemove: []string{"a.txt"},
			expectedCount: 1,
			expectedPaths: []string{"b.txt"},
		},
		{
			name: "remove last entry",
			initialEntries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
			},
			pathsToRemove: []string{"b.txt"},
			expectedCount: 1,
			expectedPaths: []string{"a.txt"},
		},
		{
			name: "remove same path multiple times (idempotent)",
			initialEntries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			pathsToRemove: []string{"file1.txt", "file1.txt"},
			expectedCount: 0,
			expectedPaths: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := IndexHeader{
				Signature:  [4]byte{'G', 'E', 'L', 'I'},
				Version:    1,
				NumEntries: uint32(len(tt.initialEntries)),
			}
			index := NewIndex(header, tt.initialEntries, "")

			for _, path := range tt.pathsToRemove {
				index.RemoveEntry(path)
			}

			if index.Header.NumEntries != uint32(tt.expectedCount) {
				t.Errorf("NumEntries = %d, want %d", index.Header.NumEntries, tt.expectedCount)
			}

			if len(index.Entries) != tt.expectedCount {
				t.Errorf("Entries length = %d, want %d", len(index.Entries), tt.expectedCount)
			}

			for i, expectedPath := range tt.expectedPaths {
				if index.Entries[i].Path != expectedPath {
					t.Errorf("Entry[%d].Path = %s, want %s", i, index.Entries[i].Path, expectedPath)
				}
			}
		})
	}
}

// TestFindEntry tests the FindEntry method
func TestFindEntry(t *testing.T) {
	tests := []struct {
		name          string
		entries       []IndexEntry
		pathToFind    string
		expectedFound bool
		expectedHash  string
	}{
		{
			name: "find existing entry",
			entries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
				{Path: "file2.txt", Hash: "def456"},
			},
			pathToFind:    "file1.txt",
			expectedFound: true,
			expectedHash:  "abc123",
		},
		{
			name: "find non-existent entry",
			entries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			pathToFind:    "nonexistent.txt",
			expectedFound: false,
		},
		{
			name:          "find in empty index",
			entries:       []IndexEntry{},
			pathToFind:    "file1.txt",
			expectedFound: false,
		},
		{
			name: "find entry in middle",
			entries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
				{Path: "c.txt", Hash: "hash_c"},
			},
			pathToFind:    "b.txt",
			expectedFound: true,
			expectedHash:  "hash_b",
		},
		{
			name: "find first entry",
			entries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
			},
			pathToFind:    "a.txt",
			expectedFound: true,
			expectedHash:  "hash_a",
		},
		{
			name: "find last entry",
			entries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
			},
			pathToFind:    "b.txt",
			expectedFound: true,
			expectedHash:  "hash_b",
		},
		{
			name: "find with empty path",
			entries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			pathToFind:    "",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := IndexHeader{
				Signature:  [4]byte{'G', 'E', 'L', 'I'},
				Version:    1,
				NumEntries: uint32(len(tt.entries)),
			}
			index := NewIndex(header, tt.entries, "")

			result := index.FindEntry(tt.pathToFind)

			if tt.expectedFound {
				if result == nil {
					t.Errorf("FindEntry(%s) returned nil, expected entry", tt.pathToFind)
					return
				}
				if result.Path != tt.pathToFind {
					t.Errorf("FindEntry(%s).Path = %s, want %s", tt.pathToFind, result.Path, tt.pathToFind)
				}
				if result.Hash != tt.expectedHash {
					t.Errorf("FindEntry(%s).Hash = %s, want %s", tt.pathToFind, result.Hash, tt.expectedHash)
				}
			} else {
				if result != nil {
					t.Errorf("FindEntry(%s) = %v, want nil", tt.pathToFind, result)
				}
			}
		})
	}
}

// TestHasEntry tests the HasEntry method
func TestHasEntry(t *testing.T) {
	tests := []struct {
		name     string
		entries  []IndexEntry
		path     string
		expected bool
	}{
		{
			name: "has existing entry",
			entries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
				{Path: "file2.txt", Hash: "def456"},
			},
			path:     "file1.txt",
			expected: true,
		},
		{
			name: "does not have non-existent entry",
			entries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			path:     "nonexistent.txt",
			expected: false,
		},
		{
			name:     "empty index",
			entries:  []IndexEntry{},
			path:     "file1.txt",
			expected: false,
		},
		{
			name: "has entry in middle",
			entries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
				{Path: "c.txt", Hash: "hash_c"},
			},
			path:     "b.txt",
			expected: true,
		},
		{
			name: "has first entry",
			entries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
			},
			path:     "a.txt",
			expected: true,
		},
		{
			name: "has last entry",
			entries: []IndexEntry{
				{Path: "a.txt", Hash: "hash_a"},
				{Path: "b.txt", Hash: "hash_b"},
			},
			path:     "b.txt",
			expected: true,
		},
		{
			name: "empty path",
			entries: []IndexEntry{
				{Path: "file1.txt", Hash: "abc123"},
			},
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := IndexHeader{
				Signature:  [4]byte{'G', 'E', 'L', 'I'},
				Version:    1,
				NumEntries: uint32(len(tt.entries)),
			}
			index := NewIndex(header, tt.entries, "")

			result := index.HasEntry(tt.path)

			if result != tt.expected {
				t.Errorf("HasEntry(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// TestIndexInvariants tests property-based invariants
func TestIndexInvariants(t *testing.T) {
	t.Run("AddEntry then RemoveEntry removes the entry", func(t *testing.T) {
		index := NewEmptyIndex()
		entry := IndexEntry{Path: "test.txt", Hash: "hash123"}

		index.AddEntry(entry)
		if !index.HasEntry("test.txt") {
			t.Error("Entry should exist after AddEntry")
		}

		index.RemoveEntry("test.txt")
		if index.HasEntry("test.txt") {
			t.Error("Entry should not exist after RemoveEntry")
		}

		if index.Header.NumEntries != 0 {
			t.Errorf("NumEntries = %d, want 0", index.Header.NumEntries)
		}
	})

	t.Run("AddOrUpdateEntry with same path does not create duplicates", func(t *testing.T) {
		index := NewEmptyIndex()
		entry1 := IndexEntry{Path: "test.txt", Hash: "hash1"}
		entry2 := IndexEntry{Path: "test.txt", Hash: "hash2"}

		index.AddOrUpdateEntry(entry1)
		index.AddOrUpdateEntry(entry2)

		if index.Header.NumEntries != 1 {
			t.Errorf("NumEntries = %d, want 1 (no duplicates)", index.Header.NumEntries)
		}

		found := index.FindEntry("test.txt")
		if found == nil {
			t.Fatal("Entry should exist")
		}
		if found.Hash != "hash2" {
			t.Errorf("Hash = %s, want hash2 (updated value)", found.Hash)
		}
	})

	t.Run("NumEntries always matches Entries length", func(t *testing.T) {
		index := NewEmptyIndex()

		// Test after various operations
		operations := []struct {
			name string
			op   func()
		}{
			{"after add", func() { index.AddEntry(IndexEntry{Path: "file1.txt", Hash: "h1"}) }},
			{"after second add", func() { index.AddEntry(IndexEntry{Path: "file2.txt", Hash: "h2"}) }},
			{"after update", func() { index.AddOrUpdateEntry(IndexEntry{Path: "file1.txt", Hash: "h1_updated"}) }},
			{"after remove", func() { index.RemoveEntry("file1.txt") }},
			{"after add again", func() { index.AddEntry(IndexEntry{Path: "file3.txt", Hash: "h3"}) }},
		}

		for _, op := range operations {
			op.op()
			if index.Header.NumEntries != uint32(len(index.Entries)) {
				t.Errorf("After %s: NumEntries(%d) != len(Entries)(%d)",
					op.name, index.Header.NumEntries, len(index.Entries))
			}
		}
	})

	t.Run("FindEntry returns pointer to entry in Entries slice", func(t *testing.T) {
		index := NewEmptyIndex()
		entry := IndexEntry{Path: "test.txt", Hash: "hash123", Size: 100}
		index.AddEntry(entry)

		found := index.FindEntry("test.txt")
		if found == nil {
			t.Fatal("FindEntry returned nil")
		}

		// Verify it's actually from the slice
		if found.Path != index.Entries[0].Path || found.Hash != index.Entries[0].Hash {
			t.Error("FindEntry did not return entry from Entries slice")
		}
	})

	t.Run("RemoveEntry on non-existent path is idempotent", func(t *testing.T) {
		index := NewEmptyIndex()
		index.AddEntry(IndexEntry{Path: "file1.txt", Hash: "h1"})

		initialCount := index.Header.NumEntries

		index.RemoveEntry("nonexistent.txt")

		if index.Header.NumEntries != initialCount {
			t.Errorf("NumEntries changed after removing non-existent entry: %d -> %d",
				initialCount, index.Header.NumEntries)
		}
	})
}
