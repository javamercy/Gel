package utilities

import (
	"os"
	"syscall"
)

type FileStatInfo struct {
	Device  uint32
	Inode   uint32
	UserId  uint32
	GroupId uint32
	Mode    uint32
	Size    uint32
}

func GetFileStatInfo(fileInfo *syscall.Stat_t) FileStatInfo {
	return FileStatInfo{
		Device:  uint32(fileInfo.Dev),
		Inode:   uint32(fileInfo.Ino),
		UserId:  fileInfo.Uid,
		GroupId: fileInfo.Gid,
		Mode:    uint32(fileInfo.Mode),
		Size:    uint32(fileInfo.Size),
	}
}

func GetFileStatFromPath(path string) (FileStatInfo, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return FileStatInfo{}, err
	}
	return GetFileStatInfo(&stat), nil
}

func GetFileStatFromFileInfo(info os.FileInfo) FileStatInfo {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return FileStatInfo{
			Size: uint32(info.Size()),
			Mode: uint32(info.Mode()),
		}
	}
	return GetFileStatInfo(stat)
}
