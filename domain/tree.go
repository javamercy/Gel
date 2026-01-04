package domain

import (
	"Gel/core/constant"
	"Gel/core/validation"
	"encoding/hex"
	"errors"
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

func NewTreeFromEntries(entries []TreeEntry) *Tree {
	var body []byte
	for _, entry := range entries {
		modeStr := entry.Mode.String()
		name := entry.Name
		hashBytes, _ := hex.DecodeString(entry.Hash)
		body = append(body, []byte(modeStr)...)
		body = append(body, constant.SpaceByte)
		body = append(body, []byte(name)...)
		body = append(body, constant.NullByte)
		body = append(body, hashBytes...)
	}
	return &Tree{
		body:    body,
		entries: entries,
	}
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
		for i < len(body) && body[i] != constant.SpaceByte {
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
		for i < len(body) && body[i] != constant.NullByte {
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
