package domain

import (
	"Gel/core/constant"
	"Gel/core/encoding"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"
)

var (
	ErrInvalidHashLength     = errors.New("hash must be 64 hexadecimal characters (32 bytes)")
	ErrIndexTooShort         = errors.New("index file is too short: minimum 12 bytes required for header")
	ErrInvalidIndexSignature = errors.New("invalid index signature: expected 'DIRC', file may be corrupted")
	ErrTruncatedEntryData    = errors.New("index file truncated: not enough data to read all entries")
	ErrIncorrectChecksumSize = errors.New("invalid index checksum: expected 32 bytes at end of file")
	ErrChecksumMismatch      = errors.New("index checksum verification failed: file may be corrupted")
	ErrHeaderDataTooShort    = errors.New("index header is incomplete: expected 12 bytes")
	ErrEntryDataTooShort     = errors.New("index entry is incomplete: minimum 74 bytes required")
	ErrPathNotNullTerminated = errors.New("index entry path is malformed: missing null terminator")
)

type IndexHeader struct {
	Signature  [4]byte
	Version    uint32
	NumEntries uint32
}

func NewIndexHeader(signature [4]byte, version uint32, numEntries uint32) IndexHeader {
	return IndexHeader{
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

func (indexEntry *IndexEntry) GetStage() uint16 {
	return (indexEntry.Flags >> 12) & 0x3
}

type Index struct {
	Header   IndexHeader
	Entries  []*IndexEntry
	Checksum string
}

func NewIndex(header IndexHeader, entries []*IndexEntry, checksum string) *Index {
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

func (index *Index) Serialize() ([]byte, error) {
	serializedHeader := index.serializeHeader()
	serializedEntries, err := index.serializeEntries()
	if err != nil {
		return nil, err
	}

	content := append(serializedHeader, serializedEntries...)
	checksum := encoding.ComputeHash(content)
	checksumBytes, _ := hex.DecodeString(checksum)

	result := make([]byte, 0, len(content)+len(checksumBytes))
	result = append(result, content...)
	result = append(result, checksumBytes...)

	return result, nil
}

func (index *Index) serializeHeader() []byte {
	serializedHeader := make([]byte, 12)
	indexHeader := index.Header

	copy(serializedHeader[0:4], indexHeader.Signature[:])
	binary.BigEndian.PutUint32(serializedHeader[4:8], indexHeader.Version)
	binary.BigEndian.PutUint32(serializedHeader[8:12], indexHeader.NumEntries)

	return serializedHeader
}

func (index *Index) serializeEntries() ([]byte, error) {
	var serializedEntries []byte
	for _, entry := range index.Entries {
		serializedEntry, err := serializeEntry(entry)
		if err != nil {
			return nil, err
		}
		serializedEntries = append(serializedEntries, serializedEntry...)
	}
	return serializedEntries, nil
}

func (index *Index) AddEntry(entry *IndexEntry) {
	index.Entries = append(index.Entries, entry)
	index.Header.NumEntries = uint32(len(index.Entries))
}

func (index *Index) AddOrUpdateEntry(entry *IndexEntry) {
	for i := range index.Entries {
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

func serializeEntry(entry *IndexEntry) ([]byte, error) {
	totalBytes := 0

	createdTime := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(createdTime, uint32(entry.CreatedTime.Unix()))

	createdTimeNanoseconds := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(createdTimeNanoseconds, uint32(entry.CreatedTime.Nanosecond()))

	updatedTime := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(updatedTime, uint32(entry.UpdatedTime.Unix()))

	updatedTimeNanoseconds := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(updatedTimeNanoseconds, uint32(entry.UpdatedTime.Nanosecond()))

	device := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(device, entry.Device)

	inode := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(inode, entry.Inode)

	mode := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(mode, entry.Mode)

	userId := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(userId, entry.UserId)

	groupId := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(groupId, entry.GroupId)

	size := make([]byte, 4)
	totalBytes += 4
	binary.BigEndian.PutUint32(size, entry.Size)

	hashBytes, err := hex.DecodeString(entry.Hash)
	if err != nil {
		return nil, err
	}

	if len(hashBytes) != constant.Sha256ByteLength {
		return nil, ErrInvalidHashLength
	}

	totalBytes += 32

	flags := make([]byte, 2)
	totalBytes += 2
	binary.BigEndian.PutUint16(flags, entry.Flags)

	path := []byte(entry.Path)
	path = append(path, 0)
	totalBytes += len(path)

	padding := (8 - (totalBytes % 8)) % 8
	path = append(path, make([]byte, padding)...)
	totalBytes += padding

	result := make([]byte, 0, totalBytes)
	result = append(result, createdTime...)
	result = append(result, createdTimeNanoseconds...)
	result = append(result, updatedTime...)
	result = append(result, updatedTimeNanoseconds...)
	result = append(result, device...)
	result = append(result, inode...)
	result = append(result, mode...)
	result = append(result, userId...)
	result = append(result, groupId...)
	result = append(result, size...)
	result = append(result, hashBytes...)
	result = append(result, flags...)
	result = append(result, path...)

	return result, nil
}

func DeserializeIndex(data []byte) (*Index, error) {
	if len(data) == 0 {
		return NewEmptyIndex(), nil
	}
	if len(data) < 12 {
		return nil, ErrIndexTooShort
	}

	index := &Index{}

	header, err := deserializeHeader(data[:12])
	if err != nil {
		return nil, err
	}

	index.Header = header

	if !bytes.Equal(header.Signature[:], []byte(constant.GelIndexSignature)) {
		return nil, ErrInvalidIndexSignature
	}

	numEntries := header.NumEntries
	offset := 12

	for i := uint32(0); i < numEntries; i++ {
		if offset >= len(data)-32 {
			return nil, ErrTruncatedEntryData
		}

		entry, bytesRead, err := deserializeEntry(data[offset:])
		if err != nil {
			return nil, err
		}
		index.AddEntry(entry)
		offset += bytesRead
	}

	if len(data)-offset != 32 {
		return nil, ErrIncorrectChecksumSize
	}

	expectedChecksumBytes := data[len(data)-32:]
	actualChecksum := encoding.ComputeHash(data[:len(data)-32])
	actualChecksumBytes, _ := hex.DecodeString(actualChecksum)

	if !bytes.Equal(expectedChecksumBytes, actualChecksumBytes) {
		return nil, ErrChecksumMismatch
	}

	index.Checksum = actualChecksum
	return index, nil
}

func deserializeHeader(data []byte) (IndexHeader, error) {
	var header IndexHeader
	if len(data) < 12 {
		return header, ErrHeaderDataTooShort
	}
	copy(header.Signature[:], data[0:4])
	header.Version = binary.BigEndian.Uint32(data[4:8])
	header.NumEntries = binary.BigEndian.Uint32(data[8:12])
	return header, nil
}

func deserializeEntry(data []byte) (*IndexEntry, int, error) {
	if len(data) < 74 {
		return nil, 0, ErrEntryDataTooShort
	}

	entry := &IndexEntry{}

	createdTimeUnix := int64(binary.BigEndian.Uint32(data[0:4]))
	createdTimeNanoseconds := int64(binary.BigEndian.Uint32(data[4:8]))
	entry.CreatedTime = time.Unix(createdTimeUnix, createdTimeNanoseconds)

	updatedTimeUnix := int64(binary.BigEndian.Uint32(data[8:12]))
	updatedTimeNanoseconds := int64(binary.BigEndian.Uint32(data[12:16]))
	entry.UpdatedTime = time.Unix(updatedTimeUnix, updatedTimeNanoseconds)

	entry.Device = binary.BigEndian.Uint32(data[16:20])
	entry.Inode = binary.BigEndian.Uint32(data[20:24])
	entry.Mode = binary.BigEndian.Uint32(data[24:28])
	entry.UserId = binary.BigEndian.Uint32(data[28:32])
	entry.GroupId = binary.BigEndian.Uint32(data[32:36])
	entry.Size = binary.BigEndian.Uint32(data[36:40])

	hashBytes := data[40:72]
	entry.Hash = hex.EncodeToString(hashBytes)

	entry.Flags = binary.BigEndian.Uint16(data[72:74])

	pathStart := 74
	pathEnd := pathStart
	for pathEnd < len(data) && data[pathEnd] != 0 {
		pathEnd++
	}

	if pathEnd >= len(data) {
		return nil, 0, ErrPathNotNullTerminated
	}

	entry.Path = string(data[pathStart:pathEnd])

	totalSize := 74 + len(entry.Path) + 1
	padding := (8 - (totalSize % 8)) % 8
	totalSize += padding

	return entry, totalSize, nil
}

func ComputeIndexFlags(path string, stage uint16) uint16 {
	pathLength := min(len(path), 0xFFF)
	flags := uint16(pathLength) | (stage << 12)
	return flags
}
