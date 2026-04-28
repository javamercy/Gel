package core

import (
	"Gel/internal/domain"
	"strings"
)

const refsHeadsPrefix = "refs/heads/"

type SymbolicRefService struct {
	refService *RefService
}

// NewSymbolicRefService creates a symbolic-ref service backed by RefService.
func NewSymbolicRefService(refService *RefService) *SymbolicRefService {
	return &SymbolicRefService{
		refService: refService,
	}
}

// Read returns the symbolic target stored in name (for example HEAD).
// When short is true and the target is a local branch ref, it returns only the branch name.
func (s *SymbolicRefService) Read(name string, short bool) (string, error) {
	ref, err := s.refService.ReadSymbolic(name)
	if err != nil {
		return "", err
	}
	if short && strings.HasPrefix(ref, refsHeadsPrefix) {
		return strings.TrimPrefix(ref, refsHeadsPrefix), nil
	}
	return ref, nil
}

// Write updates a symbolic reference file to point at ref.
// For HEAD, writes should point to refs/heads/<branch>.
func (s *SymbolicRefService) Write(name, ref string) error {
	if name == domain.HeadFileName && !strings.HasPrefix(ref, refsHeadsPrefix) {
		return ErrInvalidSymbolicRef
	}
	return s.refService.WriteSymbolic(name, ref)
}
