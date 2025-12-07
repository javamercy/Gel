package utilities

import "Gel/src/gel/core/constant"

func ConvertModeToString(mode uint32) string {
	switch mode {
	case constant.GelExecutableFileModeOctal:
		return constant.GelExecutableFileModeStr
	case constant.GelRegularFileModeOctal:
		return constant.GelRegularFileModeStr
	case constant.GelDirectoryModeOctal:
		return constant.GelDirectoryModeStr
	case constant.GelSymlinkModeOctal:
		return constant.GelSymlinkModeStr
	default:
		return ""
	}
}

func ConvertModeToUint32(mode string) uint32 {
	switch mode {
	case constant.GelExecutableFileModeStr:
		return constant.GelExecutableFileModeOctal
	case constant.GelRegularFileModeStr:
		return constant.GelRegularFileModeOctal
	case constant.GelDirectoryModeStr:
		return constant.GelDirectoryModeOctal
	case constant.GelSymlinkModeStr:
		return constant.GelSymlinkModeOctal
	default:
		return 0
	}
}

func ConvertFilesystemModeToGelMode(mode uint32) uint32 {
	if mode&0111 != 0 {
		return constant.GelExecutableFileModeOctal
	}

	return constant.GelRegularFileModeOctal
}
