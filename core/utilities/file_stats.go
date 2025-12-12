package utilities

import (
	"syscall"
	"time"
)

type FileStatInfo struct {
	Device      uint32
	Inode       uint32
	UserId      uint32
	GroupId     uint32
	Mode        uint32
	Size        uint32
	CreatedTime time.Time
	UpdatedTime time.Time
}

func GetFileStatFromPath(path string) FileStatInfo {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return FileStatInfo{
			Device:      0,
			Inode:       0,
			UserId:      0,
			GroupId:     0,
			Mode:        0,
			Size:        0,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		}
	}
	return getFileStatInfo(&stat)
}

func getFileStatInfo(fileInfo *syscall.Stat_t) FileStatInfo {
	return FileStatInfo{
		Device:      uint32(fileInfo.Dev),
		Inode:       uint32(fileInfo.Ino),
		UserId:      fileInfo.Uid,
		GroupId:     fileInfo.Gid,
		Mode:        uint32(fileInfo.Mode),
		Size:        uint32(fileInfo.Size),
		CreatedTime: time.Unix(fileInfo.Ctimespec.Sec, fileInfo.Ctimespec.Nsec),
		UpdatedTime: time.Unix(fileInfo.Mtimespec.Sec, fileInfo.Mtimespec.Nsec),
	}
}
