package domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWorkspace_FindsParentGelDir(t *testing.T) {
	repo := t.TempDir()
	gel := filepath.Join(repo, GelDirName)
	require.NoError(t, os.Mkdir(gel, 0o755))

	nested := filepath.Join(repo, "a", "b", "c")
	require.NoError(t, os.MkdirAll(nested, 0o755))

	ws, err := NewWorkspace(nested)
	require.NoError(t, err)
	assert.Equal(t, gel, ws.GelDir)
	assert.Equal(t, filepath.Join(gel, ObjectsDirName), ws.ObjectsDir)
	assert.Equal(t, filepath.Join(gel, RefsDirName), ws.RefsDir)
	assert.Equal(t, filepath.Join(gel, IndexFileName), ws.IndexPath)
	assert.Equal(t, repo, ws.RepoDir)
	assert.Equal(t, filepath.Join(gel, ConfigFileName), ws.ConfigPath)
}

func TestNewWorkspace_NotRepository(t *testing.T) {
	_, err := NewWorkspace(t.TempDir())
	assert.ErrorIs(t, err, ErrNotAGelRepository)
}

func TestFindGelDir_IgnoresNonDirectoryGel(t *testing.T) {
	repo := t.TempDir()
	notDir := filepath.Join(repo, GelDirName)
	require.NoError(t, os.WriteFile(notDir, []byte("x"), 0o644))

	_, err := findGelDir(repo)
	assert.ErrorIs(t, err, ErrNotAGelRepository)
}
