package core

import (
	"Gel/domain"
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

func (u *UpdateRefService) Update(ref string, newHash, oldHash domain.Hash) error {
	if len(oldHash[:]) == 0 {
		return u.refService.Write(ref, newHash)
	}
	return u.updateSafe(ref, newHash, oldHash)
}

func (u *UpdateRefService) updateSafe(ref string, newHash, oldHash domain.Hash) error {
	currentHash, err := u.refService.Read(ref)
	if err != nil {
		return err
	}
	if currentHash != oldHash {
		return fmt.Errorf("'%s': %w", ref, ErrRefUpdateConflict)
	}
	return u.refService.Write(ref, newHash)
}

func (u *UpdateRefService) Delete(ref string, oldHash domain.Hash) error {
	if len(oldHash[:]) == 0 {
		return u.refService.Delete(ref)
	}
	return u.deleteSafe(ref, oldHash)
}

func (u *UpdateRefService) deleteSafe(ref string, oldHash domain.Hash) error {
	currentHash, err := u.refService.Read(ref)
	if err != nil {
		return err
	}
	if currentHash != oldHash {
		return fmt.Errorf("'%s': %w", ref, ErrRefUpdateConflict)
	}
	return u.refService.Delete(ref)
}
