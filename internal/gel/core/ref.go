package core

import (
	workspace2 "Gel/internal/gel/workspace"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrRefNotFound = errors.New("reference not found")
)

type RefService struct {
	workspaceProvider *workspace2.Provider
}

func NewRefService(workspaceProvider *workspace2.Provider) *RefService {
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
			return "", fmt.Errorf("%w: %s", ErrRefNotFound, name)
		}
		return "", err
	}

	contentStr := strings.TrimSpace(string(contentBytes))
	if !strings.HasPrefix(contentStr, "ref: ") {
		return "", fmt.Errorf("invalid symbolic ref: %s", contentStr)
	}
	return strings.TrimPrefix(contentStr, "ref: "), nil
}

func (r *RefService) WriteSymbolic(name, ref string) error {
	if name == "" || ref == "" {
		return fmt.Errorf("symbolic-ref: name and ref are required")
	}

	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("symbolic-ref: ref must start with refs/")
	}

	ws := r.workspaceProvider.GetWorkspace()
	path := filepath.Join(ws.GelDir, name)
	contentStr := fmt.Sprintf("ref: %s\n", ref)
	return os.WriteFile(path, []byte(contentStr), workspace2.FilePermission)
}

func (r *RefService) Read(ref string) (string, error) {
	if !strings.HasPrefix(ref, "refs/") {
		return "", fmt.Errorf("ref must start with refs/")
	}
	ws := r.workspaceProvider.GetWorkspace()
	absPath := filepath.Join(ws.GelDir, ref)
	contentBytes, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("%w: %s", ErrRefNotFound, ref)
		}
		return "", err
	}
	if len(contentBytes) == 0 {
		return "", nil
	}
	hash := strings.TrimSpace(string(contentBytes))
	return hash, nil
}

func (r *RefService) Write(ref, hash string) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("ref must start with refs/")
	}
	ws := r.workspaceProvider.GetWorkspace()
	absPath := filepath.Join(ws.GelDir, ref)
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, workspace2.DirPermission); err != nil {
		return err
	}

	contentStr := fmt.Sprintf("%s\n", hash)
	return os.WriteFile(absPath, []byte(contentStr), workspace2.FilePermission)
}

func (r *RefService) Delete(ref string) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("ref must start with refs/")
	}

	ws := r.workspaceProvider.GetWorkspace()
	path := filepath.Join(ws.GelDir, ref)
	return os.RemoveAll(path)
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
