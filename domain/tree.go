package domain

import (
	"Gel/core/constant"
	"encoding/hex"
)

type TreeEntry struct {
	Mode string
	Hash string
	Name string
}

func NewTreeEntry(mode, hash, name string) *TreeEntry {
	return &TreeEntry{
		Mode: mode,
		Hash: hash,
		Name: name,
	}
}

type Tree struct {
	*BaseObject
}

func NewTree(data []byte) *Tree {
	return &Tree{
		BaseObject: &BaseObject{
			objectType: ObjectTypeTree,
			data:       data,
		},
	}
}

func (tree *Tree) DeserializeTree() ([]*TreeEntry, error) {
	data := tree.data
	var entries []*TreeEntry
	i := 0
	for i < len(data) {
		modeStart := i
		for data[i] != constant.SpaceByte {
			i++
		}
		mode := string(data[modeStart:i])
		i++

		nameStart := i
		for data[i] != constant.NullByte {
			i++
		}
		name := string(data[nameStart:i])
		i++

		hashBytes := data[i : i+32]
		hash := hex.EncodeToString(hashBytes)
		i += 32
		entry := NewTreeEntry(mode, hash, name)
		entries = append(entries, entry)
	}
	return entries, nil
}
