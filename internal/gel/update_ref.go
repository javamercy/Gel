package gel

import (
	"Gel/internal/gel/validate"
	"fmt"
)

type UpdateRefService struct {
	refService *RefService
}

func NewUpdateRefService(refService *RefService) *UpdateRefService {
	return &UpdateRefService{
		refService: refService,
	}
}

func (s *UpdateRefService) Update(ref string, hash string) error {
	return s.refService.Write(ref, hash)
}

func (s *UpdateRefService) UpdateSafe(ref string, newHash, oldHash string) error {
	if err := validate.Hash(newHash); err != nil {
		return err
	}
	if err := validate.Hash(oldHash); err != nil {
		return err
	}

	currHash, err := s.refService.Read(ref)
	if err != nil {
		return err
	}
	if currHash != oldHash {
		return fmt.Errorf("cannot update ref '%s' because it is not pointing to '%s'", ref, oldHash)
	}
	return s.Update(ref, newHash)
}
