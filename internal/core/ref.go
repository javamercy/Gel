package core

import (
	"Gel/internal/workspace"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RefService struct {
	workspaceProvider *workspace.Provider
}

func NewRefService(workspaceProvider *workspace.Provider) *RefService {
	return &RefService{
		workspaceProvider: workspaceProvider,
	}
}

func (r *RefService) ReadSymbolic(name string) (string, error) {
	ws := r.workspaceProvider.GetWorkspace()
	refPath := filepath.Join(ws.GelDir, name)
	contentBytes, err := os.ReadFile(refPath)
	if err != nil {
		if os.IsNotExist(err) {
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

	ws := r.workspaceProvider.GetWorkspace()
	path := filepath.Join(ws.GelDir, name)
	contentStr := fmt.Sprintf("ref: %s\n", ref)
	if err := os.WriteFile(path, []byte(contentStr), workspace.FilePermission); err != nil {
		return fmt.Errorf("ref: failed to write symbolic ref '%s': %w", name, err)
	}
	return nil
}

func (r *RefService) Read(ref string) (string, error) {
	if !strings.HasPrefix(ref, "refs/") {
		return "", fmt.Errorf("'%s': %w", ref, ErrInvalidRef)
	}
	ws := r.workspaceProvider.GetWorkspace()
	absPath := filepath.Join(ws.GelDir, ref)
	contentBytes, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("'%s': %w", ref, ErrRefNotFound)
		}
		return "", fmt.Errorf("ref: failed to read '%s': %w", ref, err)
	}
	if len(contentBytes) == 0 {
		return "", nil
	}
	hash := strings.TrimSpace(string(contentBytes))
	return hash, nil
}

func (r *RefService) Write(ref, hash string) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("'%s': %w", ref, ErrInvalidRef)
	}
	ws := r.workspaceProvider.GetWorkspace()
	absPath := filepath.Join(ws.GelDir, ref)
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
		return fmt.Errorf("ref: failed to create directory '%s': %w", dir, err)
	}

	contentStr := fmt.Sprintf("%s\n", hash)
	if err := os.WriteFile(absPath, []byte(contentStr), workspace.FilePermission); err != nil {
		return fmt.Errorf("ref: failed to write '%s': %w", ref, err)
	}
	return nil
}

func (r *RefService) Delete(ref string) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("'%s': %w", ref, ErrInvalidRef)
	}

	ws := r.workspaceProvider.GetWorkspace()
	path := filepath.Join(ws.GelDir, ref)
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("ref: failed to delete '%s': %w", ref, err)
	}
	return nil
}

func (r *RefService) Exists(ref string) bool {
	ws := r.workspaceProvider.GetWorkspace()
	path := filepath.Join(ws.GelDir, ref)
	_, err := os.Stat(path)
	return err == nil
}

func (r *RefService) Resolve(name string) (string, error) {
	ref, err := r.ReadSymbolic(name)
	if err != nil {
		return "", err
	}
	return r.Read(ref)
}
