package gelErrors

type NotRepositoryError struct {
	*GelError
	Path string
}

func NewNotRepositoryError(path, message string) *NotRepositoryError {
	return &NotRepositoryError{
		GelError: NewGelError(ExitCodeFatal, message),
		Path:     path,
	}
}

type RepositoryExistsError struct {
	*GelError
	Path string
}

func NewRepositoryExistsError(path, message string) *RepositoryExistsError {
	return &RepositoryExistsError{
		GelError: NewGelError(ExitCodeFatal, message),
		Path:     path,
	}
}
