package gelErrors

type PathNotFoundError struct {
	*GelError
	Path string
}

func NewPathNotFoundError(path, message string) *PathNotFoundError {
	return &PathNotFoundError{
		GelError: NewGelError(ExitCodeFatal, message),
		Path:     path,
	}
}

type PathOutsideRepositoryError struct {
	*GelError
	Path string
}

func NewPathOutsideRepositoryError(path, message string) *PathOutsideRepositoryError {
	return &PathOutsideRepositoryError{
		GelError: NewGelError(ExitCodeFatal, message),
		Path:     path,
	}
}
