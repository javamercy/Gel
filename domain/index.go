package domain

import (
	"Gel/core/constant"
	"Gel/core/encoding"
	"Gel/core/validation"
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

const (
	IndexHeaderSignatureSize  int = 4
	IndexHeaderVersionSize    int = 4
	IndexHeaderNumEntriesSize int = 4
	IndexHeaderSize               = IndexHeaderSignatureSize + IndexHeaderVersionSize + IndexHeaderNumEntriesSize
)

const (
	IndexChecksumSize = 32
	PaddingAlignment  = 8
)

const (
	IndexEntryTimeSize         = 4
	IndexEntryDeviceSize       = 4
	IndexEntryInodeSize        = 4
	IndexEntryModeSize         = 4
	IndexEntryUidSize          = 4
	IndexEntryGidSize          = 4
	IndexEntrySizeFieldSize    = 4
	IndexEntryHashSize         = constant.Sha256ByteLength
	IndexEntryFlagsSize        = 2
	IndexEntryPathNullTermSize = 1
	IndexEntryFixedSize        = 10*IndexEntryTimeSize + IndexEntryHashSize + IndexEntryFlagsSize
	IndexEntryHashOffset       = 10 * IndexEntryTimeSize
	IndexEntryFlagsOffset      = IndexEntryHashOffset + IndexEntryHashSize
	IndexEntryPathOffset       = IndexEntryFlagsOffset + IndexEntryFlagsSize
)

const (
	MaxPathLength = 0xFFF
	StageMask     = 0x3
	StageShift    = 12
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
	Path        string `validate:"required,relativepath"`
	Hash        string `validate:"required,sha256hex"`
	Size        uint32 `validate:"gte=0"`
	Mode        uint32 `validate:"required"`
	Device      uint32
	Inode       uint32
	UserId      uint32
	GroupId     uint32
	Flags       uint16
	CreatedTime time.Time `validate:"required"`
	UpdatedTime time.Time `validate:"required"`
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
	updatedTime time.Time) (*IndexEntry, error) {
	entry := IndexEntry{
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
	validator := validation.GetValidator()
	if err := validator.Struct(entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func (indexEntry *IndexEntry) GetStage() uint16 {
	return (indexEntry.Flags >> StageShift) & StageMask
}

func (indexEntry *IndexEntry) serialize() ([]byte, error) {
	totalBytes := 0

	createdTime := make([]byte, IndexEntryTimeSize)
	totalBytes += IndexEntryTimeSize
	binary.BigEndian.PutUint32(createdTime, uint32(indexEntry.CreatedTime.Unix()))

	createdTimeNanoseconds := make([]byte, IndexEntryTimeSize)
	totalBytes += IndexEntryTimeSize
	binary.BigEndian.PutUint32(createdTimeNanoseconds, uint32(indexEntry.CreatedTime.Nanosecond()))

	updatedTime := make([]byte, IndexEntryTimeSize)
	totalBytes += IndexEntryTimeSize
	binary.BigEndian.PutUint32(updatedTime, uint32(indexEntry.UpdatedTime.Unix()))

	updatedTimeNanoseconds := make([]byte, IndexEntryTimeSize)
	totalBytes += IndexEntryTimeSize
	binary.BigEndian.PutUint32(updatedTimeNanoseconds, uint32(indexEntry.UpdatedTime.Nanosecond()))

	device := make([]byte, IndexEntryDeviceSize)
	totalBytes += IndexEntryDeviceSize
	binary.BigEndian.PutUint32(device, indexEntry.Device)

	inode := make([]byte, IndexEntryInodeSize)
	totalBytes += IndexEntryInodeSize
	binary.BigEndian.PutUint32(inode, indexEntry.Inode)

	mode := make([]byte, IndexEntryModeSize)
	totalBytes += IndexEntryModeSize
	binary.BigEndian.PutUint32(mode, indexEntry.Mode)

	userId := make([]byte, IndexEntryUidSize)
	totalBytes += IndexEntryUidSize
	binary.BigEndian.PutUint32(userId, indexEntry.UserId)

	groupId := make([]byte, IndexEntryGidSize)
	totalBytes += IndexEntryGidSize
	binary.BigEndian.PutUint32(groupId, indexEntry.GroupId)

	size := make([]byte, IndexEntrySizeFieldSize)
	totalBytes += IndexEntrySizeFieldSize
	binary.BigEndian.PutUint32(size, indexEntry.Size)

	hashBytes, err := hex.DecodeString(indexEntry.Hash)
	if err != nil {
		return nil, err
	}
	if len(hashBytes) != IndexEntryHashSize {
		return nil, ErrInvalidHashLength
	}

	totalBytes += IndexEntryHashSize

	flags := make([]byte, IndexEntryFlagsSize)
	totalBytes += IndexEntryFlagsSize
	binary.BigEndian.PutUint16(flags, indexEntry.Flags)

	path := []byte(indexEntry.Path)
	path = append(path, 0)
	totalBytes += len(path)

	padding := (PaddingAlignment - (totalBytes % PaddingAlignment)) % PaddingAlignment
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
	checksum := encoding.ComputeSha256(content)
	checksumBytes, _ := hex.DecodeString(checksum)

	result := make([]byte, 0, len(content)+len(checksumBytes))
	result = append(result, content...)
	result = append(result, checksumBytes...)

	return result, nil
}

func (index *Index) serializeHeader() []byte {
	serializedHeader := make([]byte, IndexHeaderSize)
	header := index.Header

	copy(serializedHeader[0:IndexHeaderSignatureSize], header.Signature[:])
	binary.BigEndian.PutUint32(serializedHeader[IndexHeaderSignatureSize:IndexHeaderSignatureSize+IndexHeaderVersionSize], header.Version)
	binary.BigEndian.PutUint32(serializedHeader[IndexHeaderSignatureSize+IndexHeaderVersionSize:], header.NumEntries)

	return serializedHeader
}

func (index *Index) serializeEntries() ([]byte, error) {
	var serializedEntries []byte
	for _, entry := range index.Entries {
		serializedEntry, err := entry.serialize()
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

func DeserializeIndex(data []byte) (*Index, error) {
	if len(data) == 0 {
		return NewEmptyIndex(), nil
	}
	if len(data) < IndexHeaderSize {
		return nil, ErrIndexTooShort
	}

	var index Index

	header, err := deserializeHeader(data[:IndexHeaderSize])
	if err != nil {
		return nil, err
	}

	index.Header = header

	if !bytes.Equal(header.Signature[:], []byte(constant.GelIndexSignature)) {
		return nil, ErrInvalidIndexSignature
	}

	numEntries := header.NumEntries
	offset := IndexHeaderSize

	for i := uint32(0); i < numEntries; i++ {
		if offset >= len(data)-IndexChecksumSize {
			return nil, ErrTruncatedEntryData
		}

		entry, bytesRead, err := deserializeEntry(data[offset:])
		if err != nil {
			return nil, err
		}
		index.AddEntry(entry)
		offset += bytesRead
	}

	if len(data)-offset != IndexChecksumSize {
		return nil, ErrIncorrectChecksumSize
	}

	expectedChecksumBytes := data[len(data)-IndexChecksumSize:]
	actualChecksum := encoding.ComputeSha256(data[:len(data)-IndexChecksumSize])
	actualChecksumBytes, err := hex.DecodeString(actualChecksum)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(expectedChecksumBytes, actualChecksumBytes) {
		return nil, ErrChecksumMismatch
	}

	index.Checksum = actualChecksum
	return &index, nil
}

func deserializeHeader(data []byte) (IndexHeader, error) {
	var header IndexHeader
	if len(data) < IndexHeaderSize {
		return header, ErrHeaderDataTooShort
	}
	copy(header.Signature[:], data[0:4])
	header.Version = binary.BigEndian.Uint32(data[4:8])
	header.NumEntries = binary.BigEndian.Uint32(data[8:12])
	return header, nil
}

func deserializeEntry(data []byte) (*IndexEntry, int, error) {
	if len(data) < IndexEntryFixedSize {
		return nil, 0, ErrEntryDataTooShort
	}

	var entry IndexEntry

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

	hashBytes := data[IndexEntryHashOffset : IndexEntryHashOffset+IndexEntryHashSize]
	entry.Hash = hex.EncodeToString(hashBytes)

	entry.Flags = binary.BigEndian.Uint16(data[IndexEntryFlagsOffset : IndexEntryFlagsOffset+IndexEntryFlagsSize])

	pathStart := IndexEntryPathOffset
	pathEnd := pathStart
	for pathEnd < len(data) && data[pathEnd] != 0 {
		pathEnd++
	}

	if pathEnd >= len(data) {
		return nil, 0, ErrPathNotNullTerminated
	}

	entry.Path = string(data[pathStart:pathEnd])

	validator := validation.GetValidator()
	if err := validator.Struct(entry); err != nil {
		return nil, 0, err
	}

	totalSize := IndexEntryPathOffset + len(entry.Path) + IndexEntryPathNullTermSize
	padding := (PaddingAlignment - (totalSize % PaddingAlignment)) % PaddingAlignment
	totalSize += padding

	return &entry, totalSize, nil
}

func ComputeIndexFlags(path string, stage uint16) uint16 {
	pathLength := min(len(path), MaxPathLength)
	flags := uint16(pathLength) | (stage << StageShift)
	return flags
}
