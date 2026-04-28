package domain

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNormalizedPathUnchecked(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "valid nested", path: "src/main.go"},
		{name: "root path", path: ""},
		{name: "absolute", path: "/etc/passwd", wantErr: true},
		{name: "windows separators", path: "src\\main.go", wantErr: true},
		{name: "null byte", path: "a\x00b", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				p, err := NewNormalizedPathUnchecked(tt.path)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tt.path, p.String())
			},
		)
	}
}

func TestAbsoluteAndNormalizedConversions(t *testing.T) {
	repoDir := t.TempDir()
	filePath := filepath.Join(repoDir, "dir", "file.txt")

	abs, err := NewAbsolutePath(filePath)
	require.NoError(t, err)

	norm, err := abs.ToNormalizedPath(repoDir)
	require.NoError(t, err)
	assert.Equal(t, NormalizedPath("dir/file.txt"), norm)

	back, err := norm.ToAbsolutePath(repoDir)
	require.NoError(t, err)
	assert.Equal(t, abs.String(), back.String())
}

func TestToNormalizedPath_Root(t *testing.T) {
	repoDir := t.TempDir()
	abs, err := NewAbsolutePath(repoDir)
	require.NoError(t, err)

	norm, err := abs.ToNormalizedPath(repoDir)
	require.NoError(t, err)
	assert.Equal(t, RootPath, norm)
}

func TestNewNormalizedPath_FromRepoAndPath(t *testing.T) {
	repoDir := t.TempDir()
	input := filepath.Join(repoDir, "pkg", "x.go")

	norm, err := NewNormalizedPath(repoDir, input)
	require.NoError(t, err)
	assert.Equal(t, NormalizedPath("pkg/x.go"), norm)
}
