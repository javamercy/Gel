package gelErrors

type ObjectNotFoundError struct {
	*GelError
	Hash string
}

func NewObjectNotFoundError(hash, message string) *ObjectNotFoundError {
	return &ObjectNotFoundError{
		GelError: NewGelError(ExitCodeFatal, message),
		Hash:     hash,
	}
}

type InvalidObjectError struct {
	*GelError
	Hash string
}

func NewInvalidObjectError(hash, message string) *InvalidObjectError {
	return &InvalidObjectError{
		GelError: NewGelError(ExitCodeFatal, message),
		Hash:     hash,
	}
}
