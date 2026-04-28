package domain

import (
	"fmt"
	"syscall"
	"time"
)

// FileStat holds metadata about a file, derived from a syscall stat call.
// Used for tracking file identity and detecting changes in the index.
type FileStat struct {
	// Device is the ID of the device containing the file.
	Device uint64
	// Inode is the file's inode number.
	Inode uint64
	// UserID is the file owner's user ID.
	UserID uint32
	// GroupID is the file owner's group ID.
	GroupID uint32
	// Mode is the file mode (permissions and type).
	Mode uint32
	// Size is the file size in bytes.
	Size uint64
	// ChangedTime is the last metadata change time (ctime on Unix).
	ChangedTime time.Time
	// ModifiedTime is the last content modification time (mtime on Unix).
	ModifiedTime time.Time
}

// ParseFileStatFromPath retrieves file metadata for the given absolute path.
func ParseFileStatFromPath(path AbsolutePath) (*FileStat, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path.String(), &stat); err != nil {
		return nil, fmt.Errorf("failed to get file stat '%s': %w", path, err)
	}
	return parseFileStat(&stat), nil
}

// parseFileStat converts a syscall.Stat_t into a domain FileStat.
func parseFileStat(stat *syscall.Stat_t) *FileStat {
	return &FileStat{
		Device:       uint64(stat.Dev),
		Inode:        stat.Ino,
		UserID:       stat.Uid,
		GroupID:      stat.Gid,
		Mode:         uint32(stat.Mode),
		Size:         uint64(stat.Size),
		ChangedTime:  time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec),
		ModifiedTime: time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec),
	}
}
