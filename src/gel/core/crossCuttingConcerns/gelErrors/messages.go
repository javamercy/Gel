package gelErrors

// repository error messages
const (
	MessageNotRepository            = "not a gel repository (or any of the parent directories): %s"
	MessageRepositoryAlreadyExists  = "Reinitialized existing Gel repository in %s"
	MessageRepositoryInitialized    = "Initialized empty Gel repository in %s"
	MessageRepositoryCreationFailed = "failed to create repository at %s: %v"
)

// Index error messages
const (
	MessageIndexNotFound     = "no index found"
	MessageIndexCorrupted    = "invalid index file: %s"
	MessageIndexEmpty        = "nothing to commit (staging area empty)"
	MessageIndexChecksumFail = "index checksum mismatch"
)

// Object error messages
const (
	MessageObjectNotFound      = "object %s not found"
	MessageInvalidObject       = "invalid object %s: %s"
	MessageInvalidObjectFormat = "invalid object format: %s"
	MessageObjectTypeMismatch  = "object is %s, expected %s"
	MessageInvalidObjectType   = "invalid object type: %q (must be blob, tree, or commit)"
)

// Path error messages
const (
	MessagePathNotFound          = "pathspec '%s' did not match any files"
	MessagePathOutsideRepository = "'%s' is outside repository"
	MessagePathIsDirectory       = "'%s' is a directory"
	MessagePathAccessDenied      = "permission denied: %s"
)

// Validation error messages
const (
	MessageValidationFailed = "validation failed for %s: %s"
	MessageMissingArgument  = "missing required argument: %s"
	MessageInvalidArgument  = "invalid %s: %s"
)

// Generic messages
const (
	MessageOperationFailed = "operation failed: %v"
	MessageUnexpectedError = "unexpected error occurred: %v"
)
