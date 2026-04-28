package domain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidTreeEntryName is returned when a tree entry name is invalid
	// (empty, contains slash, null bytes, or path traversal components).
	ErrInvalidTreeEntryName = errors.New("invalid tree entry name")

	// ErrTreeMissingNullByte is returned when a tree entry name is not null-terminated.
	ErrTreeMissingNullByte = errors.New("invalid tree format: missing null byte after name")

	// ErrTreeTruncatedHash is returned when a tree entry hash is incomplete.
	ErrTreeTruncatedHash = errors.New("invalid tree format: truncated hash")
)

// TreeEntry represents a single entry within a tree object.
// Each entry corresponds to a file or subdirectory.
type TreeEntry struct {
	// Mode is the file mode (e.g., 100644 for regular file, 040000 for directory).
	Mode FileMode
	// Hash is the SHA-256 content hash of the referenced object.
	Hash Hash
	// Name is the entry's filename (not a path).
	Name string
}

// NewTreeEntry returns a TreeEntry without validating name or mode.
func NewTreeEntry(mode FileMode, hash Hash, name string) TreeEntry {
	return TreeEntry{
		Mode: mode,
		Hash: hash,
		Name: name,
	}
}

// Tree represents a directory structure in the object database.
// It contains a list of entries (files and subdirectories) with their modes and hashes.
// body stores the raw tree payload, entries caches the parsed entries.
type Tree struct {
	body    []byte
	entries []TreeEntry
}

// Body returns a defensive copy of the raw tree body bytes.
func (t *Tree) Body() []byte {
	return append([]byte(nil), t.body...)
}

// NewTree parses raw tree body bytes and returns a validated Tree.
// The input is copied to prevent external mutation.
func NewTree(body []byte) (*Tree, error) {
	bodyCopy := append([]byte(nil), body...)
	tree := &Tree{
		body: bodyCopy,
	}
	entries, err := tree.Deserialize()
	if err != nil {
		return nil, err
	}
	tree.entries = entries
	return tree, nil
}

// NewTreeFromEntries validates entries and returns a Tree built from them.
func NewTreeFromEntries(entries []TreeEntry) (*Tree, error) {
	var buffer bytes.Buffer
	for _, entry := range entries {
		if err := validateTreeEntryName(entry.Name); err != nil {
			return nil, err
		}
		if !entry.Mode.IsValid() {
			return nil, ErrInvalidFileMode
		}

		buffer.Write([]byte(fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name)))
		buffer.Write(entry.Hash[:])
	}
	return &Tree{
		body:    buffer.Bytes(),
		entries: entries,
	}, nil
}

// Type returns the domain object type for Tree.
func (t *Tree) Type() ObjectType {
	return ObjectTypeTree
}

// Size returns the byte length of the raw tree body.
func (t *Tree) Size() int {
	return len(t.body)
}

// Serialize returns the full object serialization in the form "<type> <size>\x00<body>".
func (t *Tree) Serialize() []byte {
	return SerializeObject(ObjectTypeTree, t.body)
}

// Deserialize parses the raw tree body into a slice of TreeEntry.
// It validates each entry's mode, name, and hash format.
func (t *Tree) Deserialize() ([]TreeEntry, error) {
	body := t.body
	var entries []TreeEntry
	i := 0
	for i < len(body) {
		modeStart := i
		for i < len(body) && body[i] != ' ' {
			i++
		}
		if i >= len(body) {
			return nil, ErrInvalidFileMode
		}
		modeStr := string(body[modeStart:i])
		mode := ParseFileModeFromString(modeStr)
		if !mode.IsValid() {
			return nil, ErrInvalidFileMode
		}

		i++

		nameStart := i
		for i < len(body) && body[i] != 0 {
			i++
		}
		if i >= len(body) {
			return nil, ErrTreeMissingNullByte
		}
		name := string(body[nameStart:i])
		i++

		if i+32 > len(body) {
			return nil, ErrTreeTruncatedHash
		}

		hashBytes := body[i : i+32]
		hash, err := NewHash(hex.EncodeToString(hashBytes))
		if err != nil {
			return nil, err
		}

		i += 32
		entry := NewTreeEntry(mode, hash, name)
		entries = append(entries, entry)
	}
	return entries, nil
}

// validateTreeEntryName checks that an entry name is valid for tree serialization.
// Returns ErrInvalidTreeEntryName if name is empty, contains path separators,
// null bytes, or path traversal components.
func validateTreeEntryName(name string) error {
	if name == "" || name == "." || name == ".." {
		return ErrInvalidTreeEntryName
	}
	if strings.Contains(name, "/") {
		return ErrInvalidTreeEntryName
	}
	if strings.Contains(name, "\x00") {
		return ErrInvalidTreeEntryName
	}
	return nil
}
