package utilities

import (
	"syscall"
)

type FileStatInfo struct {
	Device  uint32
	Inode   uint32
	UserId  uint32
	GroupId uint32
	Mode    uint32
}

func GetFileStatInfo(fileInfo *syscall.Stat_t) FileStatInfo {
	return FileStatInfo{
		Device:  uint32(fileInfo.Dev),
		Inode:   uint32(fileInfo.Ino),
		UserId:  fileInfo.Uid,
		GroupId: fileInfo.Gid,
		Mode:    uint32(fileInfo.Mode),
	}
}

func GetFileStatFromPath(path string) FileStatInfo {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return FileStatInfo{
			Device:  0,
			Inode:   0,
			UserId:  0,
			GroupId: 0,
			Mode:    0,
		}
	}
	return GetFileStatInfo(&stat)
}
