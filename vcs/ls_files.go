package vcs

import (
	"Gel/core/constant"
	"Gel/domain"
	"os"
	"strconv"
	"strings"
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

func (lsFilesService *LsFilesService) LsFiles(cached, stage, modified, deleted bool) (string, error) {
	index, err := lsFilesService.indexService.Read()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	entries := index.Entries

	if stage {
		lsFilesService.LsFilesWithStage(&builder, entries)
	} else if cached {
		lsFilesService.LsFilesWithCache(&builder, entries)
	} else if modified {
		lsFilesService.LsFilesWithModified(&builder, entries)
	} else if deleted {
		lsFilesService.LsFilesWithDeleted(&builder, entries)
	} else {
		lsFilesService.LsFilesWithCache(&builder, entries)
	}
	return builder.String(), nil
}

func (lsFilesService *LsFilesService) LsFilesWithStage(builder *strings.Builder, entries []*domain.IndexEntry) {
	for _, entry := range entries {
		builder.WriteString(domain.ParseFileMode(entry.Mode).String())
		builder.WriteString(constant.SpaceStr)
		builder.WriteString(entry.Hash)
		builder.WriteString(constant.SpaceStr)
		builder.WriteString(strconv.Itoa(int(entry.GetStage())))
		builder.WriteString(constant.TabStr)
		builder.WriteString(entry.Path)
		builder.WriteString(constant.NewLineStr)
	}
}

func (lsFilesService *LsFilesService) LsFilesWithCache(builder *strings.Builder, entries []*domain.IndexEntry) {
	for _, entry := range entries {
		builder.WriteString(entry.Path)
		builder.WriteString(constant.NewLineStr)
	}
}

func (lsFilesService *LsFilesService) LsFilesWithModified(builder *strings.Builder, entries []*domain.IndexEntry) {
	for _, entry := range entries {
		exists := lsFilesService.filesystemService.Exists(entry.Path)
		if !exists {
			continue
		}

		isModified := lsFilesService.isModified(entry)
		if !isModified {
			continue
		}
		builder.WriteString(entry.Path)
		builder.WriteString(constant.NewLineStr)
	}
}

func (lsFilesService *LsFilesService) LsFilesWithDeleted(builder *strings.Builder, entries []*domain.IndexEntry) {
	for _, entry := range entries {
		exists := lsFilesService.filesystemService.Exists(entry.Path)
		if !exists {
			builder.WriteString(entry.Path)
			builder.WriteString(constant.NewLineStr)
		}
	}
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
