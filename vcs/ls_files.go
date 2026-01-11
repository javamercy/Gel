package vcs

import (
	"Gel/core/constant"
	"Gel/domain"
	"io"
	"os"
	"strconv"
)

type LsFilesService struct {
	indexService      *IndexService
	filesystemService *FilesystemService
	objectService     *ObjectService
}

func NewLsFilesService(indexService *IndexService, filesystemService *FilesystemService, objectService *ObjectService) *LsFilesService {
	return &LsFilesService{
		indexService:      indexService,
		filesystemService: filesystemService,
		objectService:     objectService,
	}
}

func (lsFilesService *LsFilesService) LsFiles(w io.Writer, cached, stage, modified, deleted bool) error {
	index, err := lsFilesService.indexService.Read()
	if err != nil {
		return err
	}

	entries := index.Entries

	if stage {
		return lsFilesService.LsFilesWithStage(w, entries)
	} else if cached {
		return lsFilesService.LsFilesWithCache(w, entries)
	} else if modified {
		return lsFilesService.LsFilesWithModified(w, entries)
	} else if deleted {
		return lsFilesService.LsFilesWithDeleted(w, entries)
	}
	return lsFilesService.LsFilesWithCache(w, entries)
}

func (lsFilesService *LsFilesService) LsFilesWithStage(w io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		if _, err := io.WriteString(w, domain.ParseFileMode(entry.Mode).String()); err != nil {
			return err
		}
		if _, err := io.WriteString(w, constant.SpaceStr); err != nil {
			return err
		}
		if _, err := io.WriteString(w, entry.Hash); err != nil {
			return err
		}
		if _, err := io.WriteString(w, constant.SpaceStr); err != nil {
			return err
		}
		if _, err := io.WriteString(w, strconv.Itoa(int(entry.GetStage()))); err != nil {
			return err
		}
		if _, err := io.WriteString(w, constant.TabStr); err != nil {
			return err
		}
		if _, err := io.WriteString(w, entry.Path); err != nil {
			return err
		}
		if _, err := io.WriteString(w, constant.NewLineStr); err != nil {
			return err
		}
	}
	return nil
}

func (lsFilesService *LsFilesService) LsFilesWithCache(w io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		if _, err := io.WriteString(w, entry.Path); err != nil {
			return err
		}
		if _, err := io.WriteString(w, constant.NewLineStr); err != nil {
			return err
		}
	}
	return nil
}

func (lsFilesService *LsFilesService) LsFilesWithModified(w io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		exists := lsFilesService.filesystemService.Exists(entry.Path)
		if !exists {
			continue
		}

		isModified := lsFilesService.isModified(entry)
		if !isModified {
			continue
		}
		if _, err := io.WriteString(w, entry.Path); err != nil {
			return err
		}
		if _, err := io.WriteString(w, constant.NewLineStr); err != nil {
			return err
		}
	}
	return nil
}

func (lsFilesService *LsFilesService) LsFilesWithDeleted(w io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		exists := lsFilesService.filesystemService.Exists(entry.Path)
		if !exists {
			if _, err := io.WriteString(w, entry.Path); err != nil {
				return err
			}
			if _, err := io.WriteString(w, constant.NewLineStr); err != nil {
				return err
			}
		}
	}
	return nil
}

func (lsFilesService *LsFilesService) isModified(entry *domain.IndexEntry) bool {
	path := entry.Path
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	if uint32(stat.Size()) != entry.Size {
		return true
	}

	if !stat.ModTime().Equal(entry.UpdatedTime) {

		currentHash, err := lsFilesService.objectService.HashObject(path)
		if err != nil {
			return false
		}

		return currentHash != entry.Hash
	}
	return false
}
