package core

import (
	domain2 "Gel/internal/domain"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RefService struct {
	workspace *domain2.Workspace
}

func NewRefService(workspace *domain2.Workspace) *RefService {
	return &RefService{
		workspace: workspace,
	}
}

func (r *RefService) ReadSymbolic(name string) (string, error) {
	refPath := filepath.Join(r.workspace.GelDir, name)
	contentBytes, err := os.ReadFile(refPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("'%s': %w", name, ErrRefNotFound)
		}
		return "", fmt.Errorf("ref: failed to read '%s': %w", name, err)
	}

	contentStr := strings.TrimSpace(string(contentBytes))
	if !strings.HasPrefix(contentStr, "ref: ") {
		return "", fmt.Errorf("'%s': %w", contentStr, ErrInvalidSymbolicRef)
	}
	return strings.TrimPrefix(contentStr, "ref: "), nil
}

func (r *RefService) WriteSymbolic(name, ref string) error {
	if name == "" || ref == "" {
		return ErrInvalidSymbolicRef
	}

	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("'%s': %w", ref, ErrInvalidRef)
	}

	path := filepath.Join(r.workspace.GelDir, name)
	contentStr := fmt.Sprintf("ref: %s\n", ref)
	if err := os.WriteFile(path, []byte(contentStr), domain2.FilePermission); err != nil {
		return fmt.Errorf("ref: failed to write symbolic ref '%s': %w", name, err)
	}
	return nil
}

func (r *RefService) Read(ref string) (domain2.Hash, error) {
	if !strings.HasPrefix(ref, "refs/") {
		return domain2.Hash{}, fmt.Errorf("'%s': %w", ref, ErrInvalidRef)
	}
	absPath := filepath.Join(r.workspace.GelDir, ref)
	contentBytes, err := os.ReadFile(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain2.Hash{}, fmt.Errorf("'%s': %w", ref, ErrRefNotFound)
		}
		return domain2.Hash{}, fmt.Errorf("ref: failed to read '%s': %w", ref, err)
	}
	if len(contentBytes) == 0 {
		return domain2.Hash{}, nil
	}
	hexHash := strings.TrimSpace(string(contentBytes))
	hash, err := domain2.NewHash(hexHash)
	if err != nil {
		return domain2.Hash{}, fmt.Errorf("ref: failed to parse hash: %w", err)
	}
	return hash, nil
}

func (r *RefService) Write(ref string, hash domain2.Hash) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("'%s': %w", ref, ErrInvalidRef)
	}
	absPath := filepath.Join(r.workspace.GelDir, ref)
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, domain2.DirPermission); err != nil {
		return fmt.Errorf("ref: failed to create directory '%s': %w", dir, err)
	}

	contentStr := fmt.Sprintf("%s\n", hash)
	if err := os.WriteFile(absPath, []byte(contentStr), domain2.FilePermission); err != nil {
		return fmt.Errorf("ref: failed to write '%s': %w", ref, err)
	}
	return nil
}

func (r *RefService) Delete(ref string) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("'%s': %w", ref, ErrInvalidRef)
	}

	path := filepath.Join(r.workspace.GelDir, ref)
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("ref: failed to delete '%s': %w", ref, err)
	}
	return nil
}

func (r *RefService) Exists(ref string) bool {
	path := filepath.Join(r.workspace.GelDir, ref)
	_, err := os.Stat(path)
	return err == nil
}

func (r *RefService) Resolve(name string) (domain2.Hash, error) {
	ref, err := r.ReadSymbolic(name)
	if err != nil {
		return domain2.Hash{}, err
	}
	return r.Read(ref)
}
