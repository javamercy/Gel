package domain

import (
	"fmt"
	"syscall"
	"time"
)

// FileStat contains file system metadata used for index entries.
type FileStat struct {
	Device      uint32
	Inode       uint32
	UserId      uint32
	GroupId     uint32
	Mode        uint32
	Size        uint32
	CreatedTime time.Time
	UpdatedTime time.Time
}

// GetFileStatFromPath retrieves stat metadata for a file path.
// Errors include the path context and preserve the original cause.
func GetFileStatFromPath(path AbsolutePath) (FileStat, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path.String(), &stat); err != nil {
		return FileStat{}, fmt.Errorf("failed to get file stat '%s': %w", path, err)
	}
	return getFileStat(&stat), nil
}

// getFileStat converts a syscall.Stat_t value into the project FileStat type.
func getFileStat(stat *syscall.Stat_t) FileStat {
	return FileStat{
		Device:      uint32(stat.Dev),
		Inode:       uint32(stat.Ino),
		UserId:      stat.Uid,
		GroupId:     stat.Gid,
		Mode:        uint32(stat.Mode),
		Size:        uint32(stat.Size),
		CreatedTime: time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec),
		UpdatedTime: time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec),
	}
}
