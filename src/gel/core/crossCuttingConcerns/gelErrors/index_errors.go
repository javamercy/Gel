package gelErrors

type IndexNotFoundError struct {
	*GelError
}

func NewIndexNotFoundError(message string) *IndexNotFoundError {
	return &IndexNotFoundError{
		GelError: NewGelError(ExitCodeFatal, message),
	}
}

type IndexCorruptedError struct {
	*GelError
}

func NewIndexCorruptedError(message string) *IndexCorruptedError {
	return &IndexCorruptedError{
		GelError: NewGelError(ExitCodeFatal, message),
	}
}

type EmptyIndexError struct {
	*GelError
}

func NewEmptyIndexError(message string) *EmptyIndexError {
	return &EmptyIndexError{
		GelError: NewGelError(ExitCodeFatal, message),
	}
}
