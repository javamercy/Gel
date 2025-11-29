package dto

type AddResponse struct {
	FilesAdded []string
	Errors     []string
}

func NewAddResponse(filesAdded, errors []string) *AddResponse {
	return &AddResponse{
		FilesAdded: filesAdded,
		Errors:     errors,
	}
}
