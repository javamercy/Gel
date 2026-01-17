package domain

import (
	"Gel/core/validation"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

type TreeEntry struct {
	Mode FileMode `validate:"required"`
	Hash string   `validate:"required,sha256hex"`
	Name string   `validate:"required,relativepath"`
}

func NewTreeEntry(mode FileMode, hash, name string) (TreeEntry, error) {
	entry := TreeEntry{
		Mode: mode,
		Hash: hash,
		Name: name,
	}
	validator := validation.GetValidator()
	if err := validator.Struct(entry); err != nil {
		return TreeEntry{}, err
	}
	return entry, nil
}

type Tree struct {
	body    []byte      `validate:"required"`
	entries []TreeEntry `validate:"-"`
}

func (tree *Tree) Body() []byte {
	return tree.body
}

func NewTree(body []byte) (*Tree, error) {
	tree := &Tree{
		body: body,
	}
	validator := validation.GetValidator()
	if err := validator.Struct(tree); err != nil {
		return nil, err
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
		hashBytes, err := hex.DecodeString(entry.Hash)
		if err != nil {
			return nil, err
		}

		buffer.Write([]byte(fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name)))
		buffer.Write(hashBytes)
	}
	return &Tree{
		body:    buffer.Bytes(),
		entries: entries,
	}, nil
}

func (tree *Tree) Type() ObjectType {
	return ObjectTypeTree
}

func (tree *Tree) Size() int {
	return len(tree.body)
}

func (tree *Tree) Serialize() []byte {
	return SerializeObject(ObjectTypeTree, tree.body)
}

func (tree *Tree) Deserialize() ([]TreeEntry, error) {
	body := tree.body
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
			return nil, errors.New("invalid tree format: missing null byte after name")
		}
		name := string(body[nameStart:i])
		i++

		if i+32 > len(body) {
			return nil, errors.New("invalid tree format: truncated hash")
		}
		hashBytes := body[i : i+32]
		hash := hex.EncodeToString(hashBytes)
		i += 32
		entry, err := NewTreeEntry(mode, hash, name)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
