package services

import (
	"Gel/src/gel/core/crossCuttingConcerns/gelErrors"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/utilities"
	"Gel/src/gel/domain"
	"Gel/src/gel/domain/objects"
	"Gel/src/gel/persistence/repositories"
	"path"
	"time"
)

type IReadTreeService interface {
	ReadTree(hash string) *gelErrors.GelError
}

type ReadTreeService struct {
	objectRepository repositories.IObjectRepository
	indexRepository  repositories.IIndexRepository
}

func NewReadTreeService(objectRepository repositories.IObjectRepository, indexRepository repositories.IIndexRepository) *ReadTreeService {
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
		objectType, err := objects.GetObjectTypeByMode(treeEntry.Mode)
		fullPath := path.Join(prefix, treeEntry.Name)
		if err != nil {
			return nil, err
		}

		if objectType == objects.GelTreeObjectType {
			indexEntries, err := readTreeService.expandTree(treeEntry.Hash, fullPath)
			if err != nil {
				return nil, err
			}
			result = append(result, indexEntries...)
		} else if objectType == objects.GelBlobObjectType {

			fileStatInfo, fileStatErr := utilities.GetFileStatFromPath(fullPath)
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
				utilities.ConvertModeToUint32(treeEntry.Mode),
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

func (readTreeService *ReadTreeService) readTreeAndDeserializeTreeEntries(treeHash string) ([]*objects.TreeEntry, error) {
	compressedContent, err := readTreeService.objectRepository.Read(treeHash)
	if err != nil {
		return nil, err
	}

	content, decompressErr := encoding.Decompress(compressedContent)
	if decompressErr != nil {
		return nil, err
	}

	object, deserializeErr := objects.DeserializeObject(content)
	if deserializeErr != nil {
		return nil, err
	}

	tree, ok := object.(*objects.Tree)
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

	blob, deserializeErr := objects.DeserializeObject(content)
	if deserializeErr != nil {
		return 0, deserializeErr
	}
	return uint32(blob.Size()), nil
}
