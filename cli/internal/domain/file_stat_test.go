package domain

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFileStat(t *testing.T) {
	stat := &syscall.Stat_t{
		Dev:       42,
		Ino:       99,
		Uid:       1000,
		Gid:       1001,
		Mode:      0o100755,
		Size:      1234,
		Ctimespec: syscall.Timespec{Sec: 10, Nsec: 11},
		Mtimespec: syscall.Timespec{Sec: 20, Nsec: 21},
	}

	parsed := parseFileStat(stat)
	require.NotNil(t, parsed)
	assert.Equal(t, uint64(42), parsed.Device)
	assert.Equal(t, uint64(99), parsed.Inode)
	assert.Equal(t, uint32(1000), parsed.UserID)
	assert.Equal(t, uint32(1001), parsed.GroupID)
	assert.Equal(t, uint32(0o100755), parsed.Mode)
	assert.Equal(t, uint64(1234), parsed.Size)
	assert.Equal(t, time.Unix(10, 11), parsed.ChangedTime)
	assert.Equal(t, time.Unix(20, 21), parsed.ModifiedTime)
}

func TestParseFileStatFromPath(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "file.txt")
	require.NoError(t, os.WriteFile(p, []byte("hello"), 0o644))

	abs, err := NewAbsolutePath(p)
	require.NoError(t, err)

	fs, err := ParseFileStatFromPath(abs)
	require.NoError(t, err)
	assert.Equal(t, uint64(5), fs.Size)
	assert.NotZero(t, fs.ModifiedTime.Unix())
}

func TestParseFileStatFromPath_NotFound(t *testing.T) {
	dir := t.TempDir()
	abs, err := NewAbsolutePath(filepath.Join(dir, "missing.txt"))
	require.NoError(t, err)

	_, err = ParseFileStatFromPath(abs)
	assert.Error(t, err)
}
