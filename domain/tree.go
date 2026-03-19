package domain

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type TreeEntry struct {
	Mode FileMode
	Hash Hash
	Name string
}

func NewTreeEntry(mode FileMode, hash Hash, name string) TreeEntry {
	return TreeEntry{
		Mode: mode,
		Hash: hash,
		Name: name,
	}
}

type Tree struct {
	body    []byte
	entries []TreeEntry
}

func (t *Tree) Body() []byte {
	return t.body
}

func NewTree(body []byte) (*Tree, error) {
	tree := &Tree{
		body: body,
	}
	entries, err := tree.Deserialize()
	if err != nil {
		return nil, err
	}
	tree.entries = entries
	return tree, nil
}

func NewTreeFromEntries(entries []TreeEntry) (*Tree, error) {
	var buffer bytes.Buffer
	for _, entry := range entries {
		buffer.Write([]byte(fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name)))
		buffer.Write(entry.Hash[:])
	}
	return &Tree{
		body:    buffer.Bytes(),
		entries: entries,
	}, nil
}

func (t *Tree) Type() ObjectType {
	return ObjectTypeTree
}

func (t *Tree) Size() int {
	return len(t.body)
}

func (t *Tree) Serialize() []byte {
	return SerializeObject(ObjectTypeTree, t.body)
}

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
