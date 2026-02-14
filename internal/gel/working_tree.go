package gel

import "Gel/internal/pathspec"

type WorkingTreeService struct {
	pathResolver      *pathspec.PathResolver
	hashObjectService *HashObjectService
}

func NewWorkingDirService(
	pathResolver *pathspec.PathResolver, hashObjectService *HashObjectService,
) *WorkingTreeService {
	return &WorkingTreeService{
		pathResolver:      pathResolver,
		hashObjectService: hashObjectService,
	}
}

func (w *WorkingTreeService) GetFileMap() (map[string]string, error) {
	resolvedPaths, err := w.pathResolver.Resolve([]string{"."})
	if err != nil {
		return nil, err
	}

	results := make(map[string]string)
	for _, resolved := range resolvedPaths {
		for path := range resolved.NormalizedPaths {
			hash, _, err := w.hashObjectService.HashObject(path, false)
			if err != nil {
				return nil, err
			}
			results[path] = hash
		}
	}
	return results, nil
}
