package domain

import (
	"Gel/core/constant"
	"Gel/core/encoding"
	"Gel/domain/validation"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	ErrIndexTooShort         = errors.New("index file is too short: minimum 12 bytes required for header")
	ErrInvalidIndexSignature = errors.New("invalid index signature: expected 'DIRC', file may be corrupted")
	ErrTruncatedEntryData    = errors.New("index file truncated: not enough data to read all entries")
	ErrIncorrectChecksumSize = errors.New("invalid index checksum: expected 32 bytes at end of file")
	ErrChecksumMismatch      = errors.New("index checksum verification failed: file may be corrupted")
	ErrEntryDataTooShort     = errors.New("index entry is incomplete: minimum 74 bytes required")
	ErrPathNotNullTerminated = errors.New("index entry path is malformed: missing null terminator")
)

const (
	IndexHeaderSignatureSize  = 4
	IndexHeaderVersionSize    = 4
	IndexHeaderNumEntriesSize = 4
	IndexHeaderSize           = IndexHeaderSignatureSize + IndexHeaderVersionSize + IndexHeaderNumEntriesSize
)

const (
	IndexChecksumSize = 32
	PaddingAlignment  = 8
)

const (
	IndexEntryTimeSize              = 4
	IndexEntryDeviceSize            = 4
	IndexEntryInodeSize             = 4
	IndexEntryModeSize              = 4
	IndexEntryUserIdSize            = 4
	IndexEntryGroupIdSize           = 4
	IndexEntrySizeFieldSize         = 4
	IndexEntryHashSize              = constant.Sha256ByteLength
	IndexEntryFlagsSize             = 2
	IndexEntryPathNullTerminateSize = 1
	IndexEntryHashOffset            = 4*IndexEntryTimeSize + IndexEntryDeviceSize + IndexEntryInodeSize + IndexEntryModeSize + IndexEntryUserIdSize + IndexEntryGroupIdSize + IndexEntrySizeFieldSize
	IndexEntryFlagsOffset           = IndexEntryHashOffset + IndexEntryHashSize
	IndexEntryFixedSize             = IndexEntryFlagsOffset + IndexEntryFlagsSize
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

func NewEmptyIndexEntry(path, hash string, mode uint32) *IndexEntry {
	return &IndexEntry{
		Path:        path,
		Hash:        hash,
		Size:        0,
		Mode:        mode,
		Device:      0,
		Inode:       0,
		UserId:      0,
		GroupId:     0,
		Flags:       0,
		CreatedTime: time.Time{},
		UpdatedTime: time.Time{},
	}
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

	pathLen := len(indexEntry.Path)
	totalBytes := IndexEntryFixedSize + pathLen + IndexEntryPathNullTerminateSize
	padding := (PaddingAlignment - (totalBytes % PaddingAlignment)) % PaddingAlignment
	totalBytes += padding

	hashBytes, err := hex.DecodeString(indexEntry.Hash)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.Grow(totalBytes)
	if err := writeIndexEntryFields(&buffer, indexEntry); err != nil {
		return nil, err
	}
	if _, err := buffer.Write(hashBytes); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.BigEndian, indexEntry.Flags); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(indexEntry.Path); err != nil {
		return nil, err
	}
	buffer.WriteByte(0)
	paddingBytes := make([]byte, padding)
	if _, err := buffer.Write(paddingBytes); err != nil {
		return nil, err
	}

	serializedEntry := buffer.Bytes()
	return serializedEntry, nil
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

func (idx *Index) AddEntry(entry *IndexEntry) {
	idx.Entries = append(idx.Entries, entry)
	idx.Header.NumEntries = uint32(len(idx.Entries))
}

func (idx *Index) UpdateEntry(entry *IndexEntry) bool {
	prevEntry, i := idx.FindEntry(entry.Path)
	if prevEntry == nil {
		return false
	}
	idx.Entries[i] = entry
	return true
}

func (idx *Index) SetEntry(entry *IndexEntry) {
	if !idx.UpdateEntry(entry) {
		idx.AddEntry(entry)
	}
}

func (idx *Index) RemoveEntry(path string) {
	if entry, i := idx.FindEntry(path); entry != nil {
		idx.Entries = append(idx.Entries[:i], idx.Entries[i+1:]...)
		idx.Header.NumEntries = uint32(len(idx.Entries))
	}
}

func (idx *Index) FindEntry(path string) (*IndexEntry, int) {
	i := sort.Search(len(idx.Entries), func(i int) bool {
		return idx.Entries[i].Path >= path
	})
	if i < len(idx.Entries) && idx.Entries[i].Path == path {
		return idx.Entries[i], i
	}
	return nil, 0
}

func (idx *Index) FindEntriesByPathPrefix(prefix string) []*IndexEntry {
	var result []*IndexEntry
	for _, entry := range idx.Entries {
		if strings.HasPrefix(entry.Path, prefix) {
			result = append(result, entry)
		}
	}
	return result
}

func (idx *Index) FindEntriesByPathPattern(pattern string) []*IndexEntry {
	var result []*IndexEntry
	for _, entry := range idx.Entries {
		if match, _ := filepath.Match(pattern, entry.Path); match {
			result = append(result, entry)
		}
	}
	return result
}

func (idx *Index) HasEntry(path string) bool {
	entry, _ := idx.FindEntry(path)
	return entry != nil
}

func (idx *Index) Serialize() ([]byte, error) {
	serializedHeader := idx.serializeHeader()

	sort.Slice(idx.Entries, func(i, j int) bool {
		return idx.Entries[i].Path < idx.Entries[j].Path
	})
	serializedEntries, err := idx.serializeEntries()
	if err != nil {
		return nil, err
	}

	data := append(serializedHeader, serializedEntries...)
	checksum := encoding.ComputeSha256(data)
	checksumBytes, err := hex.DecodeString(checksum)
	if err != nil {
		return nil, err
	}

	result := make([]byte, 0, len(data)+len(checksumBytes))
	result = append(result, data...)
	result = append(result, checksumBytes...)

	return result, nil
}

func (idx *Index) serializeHeader() []byte {
	serializedHeader := make([]byte, IndexHeaderSize)
	header := idx.Header

	copy(serializedHeader[0:IndexHeaderSignatureSize], header.Signature[:])
	binary.BigEndian.PutUint32(serializedHeader[IndexHeaderSignatureSize:IndexHeaderSignatureSize+IndexHeaderVersionSize], header.Version)
	binary.BigEndian.PutUint32(serializedHeader[IndexHeaderSignatureSize+IndexHeaderVersionSize:], header.NumEntries)

	return serializedHeader
}

func (idx *Index) serializeEntries() ([]byte, error) {
	var serializedEntries []byte
	for _, entry := range idx.Entries {
		serializedEntry, err := entry.serialize()
		if err != nil {
			return nil, err
		}
		serializedEntries = append(serializedEntries, serializedEntry...)
	}
	return serializedEntries, nil
}

func DeserializeIndex(data []byte) (*Index, error) {
	if len(data) == 0 {
		return NewEmptyIndex(), nil
	}
	if len(data) < IndexHeaderSize {
		return nil, ErrIndexTooShort
	}

	var index Index
	offset := 0

	header, err := deserializeHeader(data[:IndexHeaderSize])
	if err != nil {
		return nil, err
	}
	index.Header = header

	numEntries := header.NumEntries
	offset += IndexHeaderSize

	for i := uint32(0); i < numEntries; i++ {
		if offset >= len(data)-IndexChecksumSize {
			return nil, ErrTruncatedEntryData
		}

		entry, entrySize, err := deserializeIndexEntry(data[offset:])
		if err != nil {
			return nil, err
		}

		index.AddEntry(entry)
		offset += entrySize
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

func ComputeIndexFlags(path string, stage uint16) uint16 {
	pathLength := min(len(path), MaxPathLength)
	flags := uint16(pathLength) | (stage << StageShift)
	return flags
}

func deserializeHeader(data []byte) (IndexHeader, error) {
	var header IndexHeader

	copy(header.Signature[:], data[0:IndexHeaderSignatureSize])
	if !bytes.Equal(header.Signature[:], []byte(constant.GelIndexSignature)) {
		return header, ErrInvalidIndexSignature
	}

	header.Version = binary.BigEndian.Uint32(data[IndexHeaderSignatureSize : IndexHeaderSignatureSize+IndexHeaderVersionSize])
	header.NumEntries = binary.BigEndian.Uint32(data[IndexHeaderSignatureSize+IndexHeaderVersionSize:])

	return header, nil
}

func deserializeIndexEntry(data []byte) (*IndexEntry, int, error) {
	if len(data) < IndexEntryFixedSize {
		return nil, 0, ErrEntryDataTooShort
	}

	var entry IndexEntry
	reader := bytes.NewReader(data)

	if err := readIndexEntryFields(reader, &entry); err != nil {
		return nil, 0, err
	}

	hashBytes := make([]byte, IndexEntryHashSize)
	if _, err := reader.Read(hashBytes); err != nil {
		return nil, 0, err
	}
	entry.Hash = hex.EncodeToString(hashBytes)

	if err := binary.Read(reader, binary.BigEndian, &entry.Flags); err != nil {
		return nil, 0, err
	}

	offset := IndexEntryFixedSize
	pathEnd := bytes.IndexByte(data[offset:], 0)
	if pathEnd == -1 {
		return nil, 0, ErrPathNotNullTerminated
	}

	entry.Path = string(data[offset : offset+pathEnd])
	offset += pathEnd + IndexEntryPathNullTerminateSize

	validator := validation.GetValidator()
	if err := validator.Struct(entry); err != nil {
		return nil, 0, err
	}

	padding := (PaddingAlignment - (offset % PaddingAlignment)) % PaddingAlignment
	totalSize := offset + padding

	return &entry, totalSize, nil
}

func writeIndexEntryFields(buffer *bytes.Buffer, entry *IndexEntry) error {
	if err := binary.Write(buffer, binary.BigEndian, uint32(entry.CreatedTime.Unix())); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, uint32(entry.CreatedTime.Nanosecond())); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, uint32(entry.UpdatedTime.Unix())); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, uint32(entry.UpdatedTime.Nanosecond())); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, entry.Device); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, entry.Inode); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, entry.Mode); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, entry.UserId); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, entry.GroupId); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, entry.Size); err != nil {
		return err
	}
	return nil
}

func readIndexEntryFields(reader *bytes.Reader, entry *IndexEntry) error {
	var createdTimeUnix, createdTimeNanoseconds uint32
	if err := binary.Read(reader, binary.BigEndian, &createdTimeUnix); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.BigEndian, &createdTimeNanoseconds); err != nil {
		return err
	}
	entry.CreatedTime = time.Unix(int64(createdTimeUnix), int64(createdTimeNanoseconds))

	var updatedTimeUnix, updatedTimeNanoseconds uint32
	if err := binary.Read(reader, binary.BigEndian, &updatedTimeUnix); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.BigEndian, &updatedTimeNanoseconds); err != nil {
		return err
	}
	entry.UpdatedTime = time.Unix(int64(updatedTimeUnix), int64(updatedTimeNanoseconds))

	if err := binary.Read(reader, binary.BigEndian, &entry.Device); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.BigEndian, &entry.Inode); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.BigEndian, &entry.Mode); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.BigEndian, &entry.UserId); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.BigEndian, &entry.GroupId); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.BigEndian, &entry.Size); err != nil {
		return err
	}
	return nil
}
