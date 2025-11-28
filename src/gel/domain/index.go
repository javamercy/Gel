package domain

import (
	"Gel/src/gel/core/constant"
	"time"
)

type IndexHeader struct {
	Signature  [4]byte
	Version    uint32
	NumEntries uint32
}

func NewIndexHeader(signature [4]byte, version uint32, numEntries uint32) *IndexHeader {
	return &IndexHeader{
		Signature:  signature,
		Version:    version,
		NumEntries: numEntries,
	}
}

type IndexEntry struct {
	Path        string
	Hash        string
	Size        uint32
	Mode        uint32
	Device      uint32
	Inode       uint32
	UserId      uint32
	GroupId     uint32
	Flags       uint16
	CreatedTime time.Time
	UpdatedTime time.Time
}

func NewIndexEntry(
	path string,
	hash string,
	size uint32,
	mode uint32,
	device uint32,
	inode uint32,
	userId uint32,
	groupId uint32,
	flags uint16,
	createdTime time.Time,
	updatedTime time.Time) *IndexEntry {
	return &IndexEntry{
		path,
		hash,
		size,
		mode,
		device,
		inode,
		userId,
		groupId,
		flags,
		createdTime,
		updatedTime,
	}
}

type Index struct {
	Header   *IndexHeader
	Entries  []*IndexEntry
	Checksum string
}

func NewIndex(header *IndexHeader, entries []*IndexEntry, checksum string) *Index {
	return &Index{
		Header:   header,
		Entries:  entries,
		Checksum: checksum,
	}
}

func NewEmptyIndex() *Index {
	signatureBytes := [4]byte([]byte(constant.GelIndexSignature))
	header := NewIndexHeader(signatureBytes, constant.GelIndexVersion, 0)
	return NewIndex(header, []*IndexEntry{}, "")
}

func (index *Index) AddEntry(entry *IndexEntry) {
	index.Entries = append(index.Entries, entry)
	index.Header.NumEntries = uint32(len(index.Entries))
}

func (index *Index) AddOrUpdateEntry(entry *IndexEntry) {
	for i, _ := range index.Entries {
		if index.Entries[i].Path == entry.Path {
			index.Entries[i] = entry
			return
		}
	}
	index.AddEntry(entry)
}

func (index *Index) RemoveEntry(path string) {
	for i, entry := range index.Entries {
		if entry.Path == path {
			index.Entries = append(index.Entries[:i], index.Entries[i+1:]...)
			index.Header.NumEntries = uint32(len(index.Entries))
			return
		}
	}
}

func (index *Index) FindEntry(path string) *IndexEntry {
	for _, entry := range index.Entries {
		if entry.Path == path {
			return entry
		}
	}
	return nil
}

func (index *Index) HasEntry(path string) bool {
	for _, entry := range index.Entries {
		if entry.Path == path {
			return true
		}
	}
	return false
}
