package dto

type InitRequest struct {
	Path string
}

func NewInitRequest(path string) *InitRequest {
	return &InitRequest{
		path,
	}
}
