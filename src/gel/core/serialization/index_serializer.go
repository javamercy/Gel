package serialization

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/domain"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"
)

func SerializeIndex(index *domain.Index) []byte {
	serializedHeader := serializeIndexHeader(index.Header)
	var serializedEntries []byte
	for _, entry := range index.Entries {
		serializedEntry := serializeIndexEntry(entry)
		serializedEntries = append(serializedEntries, serializedEntry...)
	}

	content := append(serializedHeader, serializedEntries...)
	checksum := encoding.ComputeHash(content)
	checksumBytes, _ := hex.DecodeString(checksum)

	result := make([]byte, 0, len(content)+len(checksumBytes))
	result = append(result, content...)
	result = append(result, checksumBytes...)

	return result
}

func serializeIndexHeader(indexHeader *domain.IndexHeader) []byte {
	header := make([]byte, 12)
	copy(header[0:4], indexHeader.Signature[:])
	binary.BigEndian.PutUint32(header[4:8], indexHeader.Version)
	binary.BigEndian.PutUint32(header[8:12], indexHeader.NumEntries)

	return header
}

func serializeIndexEntry(entry *domain.IndexEntry) []byte {
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
	if err != nil || len(hashBytes) != 32 {
		hashBytes = make([]byte, 32)
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

	return result
}

func DeserializeIndex(data []byte) (*domain.Index, error) {
	if len(data) == 0 {
		return domain.NewEmptyIndex(), nil
	}
	if len(data) < 12 {
		return nil, errors.New("invalid index file: too short for header")
	}

	index := &domain.Index{}

	header, err := deserializeIndexHeader(data[:12])
	if err != nil {
		return nil, err
	}

	index.Header = header

	if !bytes.Equal(header.Signature[:], []byte(constant.GelIndexSignature)) {
		return nil, errors.New("invalid index signature")
	}

	numEntries := header.NumEntries
	offset := 12

	for i := uint32(0); i < numEntries; i++ {
		if offset >= len(data)-32 {
			return nil, errors.New("invalid index: truncated entry data")
		}

		entry, bytesRead, err := deserializeIndexEntry(data[offset:])
		if err != nil {
			return nil, err
		}
		index.AddEntry(entry)
		offset += bytesRead
	}

	if len(data)-offset != 32 {
		return nil, errors.New("invalid index: incorrect checksum size")
	}

	expectedChecksumBytes := data[len(data)-32:]
	actualChecksum := encoding.ComputeHash(data[:len(data)-32])
	actualChecksumBytes, _ := hex.DecodeString(actualChecksum)

	if !bytes.Equal(expectedChecksumBytes, actualChecksumBytes) {
		return nil, errors.New("index checksum mismatch")
	}

	index.Checksum = actualChecksum
	return index, nil
}

func deserializeIndexHeader(data []byte) (*domain.IndexHeader, error) {
	if len(data) < 12 {
		return nil, errors.New("header data too short")
	}
	header := &domain.IndexHeader{}
	copy(header.Signature[:], data[0:4])
	header.Version = binary.BigEndian.Uint32(data[4:8])
	header.NumEntries = binary.BigEndian.Uint32(data[8:12])
	return header, nil
}

func deserializeIndexEntry(data []byte) (*domain.IndexEntry, int, error) {
	if len(data) < 74 {
		return nil, 0, errors.New("entry data too short")
	}

	entry := &domain.IndexEntry{}

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
		return nil, 0, errors.New("path not null-terminated")
	}

	entry.Path = string(data[pathStart:pathEnd])

	totalSize := 74 + len(entry.Path) + 1
	padding := (8 - (totalSize % 8)) % 8
	totalSize += padding

	return entry, totalSize, nil
}
