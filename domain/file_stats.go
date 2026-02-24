package domain

import (
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

// GetFileStatFromPath retrieves file system metadata for the given path.
func GetFileStatFromPath(path string) FileStat {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return FileStat{
			Device:      0,
			Inode:       0,
			UserId:      0,
			GroupId:     0,
			Mode:        0,
			Size:        0,
			CreatedTime: time.Time{},
			UpdatedTime: time.Time{},
		}
	}
	return getFileStat(&stat)
}

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
