package domain

import (
	"Gel/core/constant"
	"encoding/hex"
)

type TreeEntry struct {
	Mode FileMode
	Hash string
	Name string
}

func NewTreeEntry(mode FileMode, hash, name string) *TreeEntry {
	return &TreeEntry{
		Mode: mode,
		Hash: hash,
		Name: name,
	}
}

type Tree struct {
	body    []byte
	entries []*TreeEntry
}

func (tree *Tree) Body() []byte {
	return tree.body
}

func NewTree(body []byte) *Tree {
	return &Tree{
		body: body,
	}
}

func NewTreeFromEntries(entries []*TreeEntry) *Tree {
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

func (tree *Tree) Deserialize() ([]*TreeEntry, error) {
	body := tree.body
	var entries []*TreeEntry
	i := 0
	for i < len(body) {
		modeStart := i
		for body[i] != constant.SpaceByte {
			i++
		}
		modeStr := string(body[modeStart:i])
		mode := ParseFileModeFromString(modeStr)
		if !mode.IsValid() {
			return nil, ErrInvalidFileMode
		}

		i++

		nameStart := i
		for body[i] != constant.NullByte {
			i++
		}
		name := string(body[nameStart:i])
		i++

		hashBytes := body[i : i+32]
		hash := hex.EncodeToString(hashBytes)
		i += 32
		entry := NewTreeEntry(mode, hash, name)
		entries = append(entries, entry)
	}
	return entries, nil
}
