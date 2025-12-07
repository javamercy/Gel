package services

import (
	"Gel/core/crossCuttingConcerns/gelErrors"
	"Gel/core/encoding"
	utilities2 "Gel/core/utilities"
	"Gel/domain"
	objects2 "Gel/domain/objects"
	repositories2 "Gel/persistence/repositories"
	"path"
	"time"
)

type IReadTreeService interface {
	ReadTree(hash string) *gelErrors.GelError
}

type ReadTreeService struct {
	objectRepository repositories2.IObjectRepository
	indexRepository  repositories2.IIndexRepository
}

func NewReadTreeService(objectRepository repositories2.IObjectRepository, indexRepository repositories2.IIndexRepository) *ReadTreeService {
	return &ReadTreeService{
		objectRepository: objectRepository,
		indexRepository:  indexRepository,
	}
}

func (readTreeService *ReadTreeService) ReadTree(hash string) *gelErrors.GelError {

	indexEntries, expandErr := readTreeService.expandTree(hash, "")
	if expandErr != nil {
		return gelErrors.NewGelError(gelErrors.ExitCodeFatal, expandErr.Error())
	}

	index := domain.NewEmptyIndex()
	for _, entry := range indexEntries {
		index.AddEntry(entry)
	}

	writeErr := readTreeService.indexRepository.Write(index)

	if writeErr != nil {
		return gelErrors.NewGelError(gelErrors.ExitCodeFatal, writeErr.Error())
	}

	return nil

}
func (readTreeService *ReadTreeService) expandTree(treeHash, prefix string) ([]*domain.IndexEntry, error) {

	result := make([]*domain.IndexEntry, 0)
	treeEntries, gelError := readTreeService.readTreeAndDeserializeTreeEntries(treeHash)
	if gelError != nil {
		return nil, gelError
	}

	for _, treeEntry := range treeEntries {
		objectType, err := objects2.GetObjectTypeByMode(treeEntry.Mode)
		fullPath := path.Join(prefix, treeEntry.Name)
		if err != nil {
			return nil, err
		}

		if objectType == objects2.GelTreeObjectType {
			indexEntries, err := readTreeService.expandTree(treeEntry.Hash, fullPath)
			if err != nil {
				return nil, err
			}
			result = append(result, indexEntries...)
		} else if objectType == objects2.GelBlobObjectType {

			fileStatInfo, fileStatErr := utilities2.GetFileStatFromPath(fullPath)
			if fileStatErr != nil {
				return nil, fileStatErr
			}

			size, sizeErr := readTreeService.readBlobAndGetSize(treeEntry.Hash)
			if sizeErr != nil {
				return nil, sizeErr
			}

			indexEntry := domain.NewIndexEntry(
				fullPath,
				treeEntry.Hash,
				size,
				utilities2.ConvertModeToUint32(treeEntry.Mode),
				fileStatInfo.Device,
				fileStatInfo.Inode,
				fileStatInfo.UserId,
				fileStatInfo.GroupId,
				domain.ComputeIndexFlags(fullPath, 0),
				time.Now(),
				time.Now())

			result = append(result, indexEntry)
		}

	}

	return result, nil
}

func (readTreeService *ReadTreeService) readTreeAndDeserializeTreeEntries(treeHash string) ([]*objects2.TreeEntry, error) {
	compressedContent, err := readTreeService.objectRepository.Read(treeHash)
	if err != nil {
		return nil, err
	}

	content, decompressErr := encoding.Decompress(compressedContent)
	if decompressErr != nil {
		return nil, err
	}

	object, deserializeErr := objects2.DeserializeObject(content)
	if deserializeErr != nil {
		return nil, err
	}

	tree, ok := object.(*objects2.Tree)
	if !ok {
		return nil, err
	}

	treeEntries, err := tree.DeserializeTree()
	if err != nil {
		return nil, err
	}

	return treeEntries, nil
}

func (readTreeService *ReadTreeService) readBlobAndGetSize(hash string) (uint32, error) {

	compressedContent, readErr := readTreeService.objectRepository.Read(hash)
	if readErr != nil {
		return 0, readErr
	}

	content, decompressErr := encoding.Decompress(compressedContent)
	if decompressErr != nil {
		return 0, nil
	}

	blob, deserializeErr := objects2.DeserializeObject(content)
	if deserializeErr != nil {
		return 0, deserializeErr
	}
	return uint32(blob.Size()), nil
}
