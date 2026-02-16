package gel

import (
	"Gel/internal/gel/core"
	"fmt"
)

type UpdateRefService struct {
	refService *core.RefService
}

func NewUpdateRefService(refService *core.RefService) *UpdateRefService {
	return &UpdateRefService{
		refService: refService,
	}
}

func (u *UpdateRefService) Update(ref string, newHash, oldHash string) error {
	if oldHash == "" {
		return u.refService.Write(ref, newHash)
	}
	return u.updateSafe(ref, newHash, oldHash)
}

func (u *UpdateRefService) updateSafe(ref string, newHash, oldHash string) error {
	currentHash, err := u.refService.Read(ref)
	if err != nil {
		return err
	}
	if currentHash != oldHash {
		return fmt.Errorf("cannot update ref '%s' because it is not pointing to '%s'", ref, oldHash)
	}
	return u.refService.Write(ref, newHash)
}

func (u *UpdateRefService) Delete(ref string, oldHash string) error {
	if oldHash == "" {
		return u.refService.Delete(ref)
	}
	return u.deleteSafe(ref, oldHash)
}

func (u *UpdateRefService) deleteSafe(ref, oldHash string) error {
	currentHash, err := u.refService.Read(ref)
	if err != nil {
		return err
	}
	if currentHash != oldHash {
		return fmt.Errorf("cannot update ref '%s' because it is not pointing to '%s'", ref, oldHash)
	}
	return u.refService.Delete(ref)
}
