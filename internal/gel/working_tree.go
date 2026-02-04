package gel

import "Gel/internal/pathspec"

type WorkingDirService struct {
	pathResolver      *pathspec.PathResolver
	hashObjectService *HashObjectService
}

func NewWorkingDirService(
	pathResolver *pathspec.PathResolver, hashObjectService *HashObjectService,
) *WorkingDirService {
	return &WorkingDirService{
		pathResolver:      pathResolver,
		hashObjectService: hashObjectService,
	}
}

func (w *WorkingDirService) GetWorkingDirFiles() (map[string]string, error) {
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
