package services

import (
	"Gel/src/gel/application/rules"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/domain"
	"Gel/src/gel/persistence/repositories"
	"os"
)

type UpdateIndexOptions struct {
	Add    bool
	Remove bool
}

type IUpdateIndexService interface {
	UpdateIndex(paths []string, options UpdateIndexOptions) error
	add(index *domain.Index, paths []string, indexFilePath string) error
	remove(paths []string) error
}

type UpdateIndexService struct {
	gelRepository        repositories.IGelRepository
	filesystemRepository repositories.IFilesystemRepository
	updateIndexRules     *rules.UpdateIndexRules
}

func NewUpdateIndexService(gelRepository repositories.IGelRepository, filesystemRepository repositories.IFilesystemRepository, updateIndexRules *rules.UpdateIndexRules) *UpdateIndexService {
	return &UpdateIndexService{
		gelRepository,
		filesystemRepository,
		updateIndexRules,
	}
}

func (updateIndexService *UpdateIndexService) UpdateIndex(paths []string, options UpdateIndexOptions) error {

	err := updateIndexService.updateIndexRules.AllPathsMustExist(paths)
	if err != nil {
		return err
	}

	err = updateIndexService.updateIndexRules.NoDuplicatePaths(paths)
	if err != nil {
		return err
	}

	err = updateIndexService.updateIndexRules.PathsMustBeFiles(paths)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	indexFilePath, err := updateIndexService.gelRepository.FindIndexFilePath(cwd)
	indexFileExists := updateIndexService.filesystemRepository.Exists(indexFilePath)

	var indexBytes []byte
	if indexFileExists {
		indexBytes, err = updateIndexService.filesystemRepository.ReadFile(indexFilePath)
		if err != nil {
			return err
		}
	}

	index, err := serialization.DeserializeIndex(indexBytes)
	if err != nil {
		return err
	}

	if options.Add {
		err := updateIndexService.add(index, paths, indexFilePath)
		if err != nil {
			return err
		}
	} else if options.Remove {
		err := updateIndexService.remove(paths)
		if err != nil {
			return err
		}
	}
	return nil
}

func (updateIndexService *UpdateIndexService) add(index *domain.Index, paths []string, indexFilePath string) error {

	for _, path := range paths {
		fileInfo, _ := updateIndexService.filesystemRepository.Stat(path)
		data, err := updateIndexService.filesystemRepository.ReadFile(path)
		if err != nil {
			return err
		}

		exists := false
		content := serialization.SerializeObject(constant.Blob, data)
		for i := range index.Entries {
			if index.Entries[i].Path == path {
				index.Entries[i].Size = uint32(fileInfo.Size())
				index.Entries[i].Mode = uint32(fileInfo.Mode())
				index.Entries[i].UpdatedTime = fileInfo.ModTime()
				index.Entries[i].Hash = encoding.ComputeHash(content)
				exists = true
				break
			}
		}

		if exists {
			continue
		}
		// Add new entry
		newEntry := domain.IndexEntry{
			Path:        path,
			Hash:        encoding.ComputeHash(content),
			Size:        uint32(fileInfo.Size()),
			Mode:        uint32(fileInfo.Mode()),
			Device:      0,
			Inode:       0,
			UserId:      0,
			GroupId:     0,
			Flags:       0,
			CreatedTime: fileInfo.ModTime(),
			UpdatedTime: fileInfo.ModTime(),
		}

		index.Entries = append(index.Entries, newEntry)
		index.Header.NumEntries = uint32(len(index.Entries))
	}

	indexBytes := serialization.SerializeIndex(index)
	index.Checksum = encoding.ComputeHash(indexBytes)

	serializedIndex := serialization.SerializeIndex(index)
	err := updateIndexService.filesystemRepository.WriteFile(
		indexFilePath,
		serializedIndex,
		false,
		constant.FilePermission)

	return err
}

func (updateIndexService *UpdateIndexService) remove(paths []string) error {

	return nil
}
