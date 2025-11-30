package dto

type AddResponse struct {
	Paths  []string
	Errors error
}

func NewAddResponse(paths []string, err error) *AddResponse {
	return &AddResponse{
		Paths:  paths,
		Errors: err,
	}
}
