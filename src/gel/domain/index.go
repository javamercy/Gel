package domain

import "time"

type IndexHeader struct {
	Signature  [4]byte
	Version    uint32
	NumEntries uint32
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

type Index struct {
	Header   IndexHeader
	Entries  []IndexEntry
	Checksum string
}

func NewIndex(header IndexHeader, entries []IndexEntry, checksum string) *Index {
	return &Index{
		Header:   header,
		Entries:  entries,
		Checksum: checksum,
	}
}

func NewEmptyIndex() *Index {
	header := IndexHeader{
		Signature:  [4]byte{'G', 'E', 'L', 'I'},
		Version:    1,
		NumEntries: 0,
	}
	return NewIndex(header, []IndexEntry{}, "")
}

func (index *Index) AddEntry(entry IndexEntry) {
	index.Entries = append(index.Entries, entry)
	index.Header.NumEntries = uint32(len(index.Entries))
}

func (index *Index) AddOrUpdateEntry(entry IndexEntry) {
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
			return &entry
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
