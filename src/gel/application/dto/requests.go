package dto

import "Gel/src/gel/core/constant"

type InitRequest struct {
	Path string
}

func NewInitRequest(path string) *InitRequest {
	return &InitRequest{
		path,
	}
}

type HashObjectRequest struct {
	Paths      []string
	ObjectType constant.ObjectType
	Write      bool
}

func NewHashObjectRequest(paths []string, objectType constant.ObjectType, write bool) *HashObjectRequest {
	return &HashObjectRequest{
		Paths:      paths,
		ObjectType: objectType,
		Write:      write,
	}
}

type CatFileRequest struct {
	Hash      string
	ShowType  bool
	ShowSize  bool
	Pretty    bool
	CheckOnly bool
}

func NewCatFileRequest(hash string, showType, showSize, pretty, checkOnly bool) *CatFileRequest {
	return &CatFileRequest{
		Hash:      hash,
		ShowType:  showType,
		ShowSize:  showSize,
		Pretty:    pretty,
		CheckOnly: checkOnly,
	}
}

type UpdateIndexRequest struct {
	Paths  []string
	Add    bool
	Remove bool
}

func NewUpdateIndexRequest(paths []string, add, remove bool) *UpdateIndexRequest {
	{
		return &UpdateIndexRequest{
			Paths:  paths,
			Add:    add,
			Remove: remove,
		}
	}
}

type LsFilesRequest struct {
	Cached   bool
	Stage    bool
	Deleted  bool
	Modified bool
}

func NewLsFilesRequest(cached, stage, deleted, modified bool) *LsFilesRequest {
	return &LsFilesRequest{
		Cached:   cached,
		Stage:    stage,
		Deleted:  deleted,
		Modified: modified,
	}
}

type AddRequest struct {
	Pathspecs []string
	DryRun    bool
	Verbose   bool
}

func NewAddRequest(pathspecs []string, dryRun, verbose bool) *AddRequest {
	{
		return &AddRequest{
			Pathspecs: pathspecs,
			DryRun:    dryRun,
			Verbose:   verbose,
		}
	}
}
