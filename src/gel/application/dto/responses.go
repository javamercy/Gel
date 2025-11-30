package dto

type AddResponse struct {
	Paths []string
	Error error
}

func NewAddResponse(paths []string, err error) *AddResponse {
	return &AddResponse{
		Paths: paths,
		Error: err,
	}
}
