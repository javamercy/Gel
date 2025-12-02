package gelErrors

type ExitCode int

const (
	ExitCodeFatal   ExitCode = 128
	ExitCodeWarning ExitCode = 0
	ExitCodeUsage   ExitCode = 129
	ExitCodeGeneral ExitCode = 1
)

type GelError struct {
	ExitCode ExitCode
	Message  string
}

func NewGelError(code ExitCode, message string) *GelError {
	return &GelError{
		ExitCode: code,
		Message:  message,
	}
}

func (gelError *GelError) GetExitCode() int {
	return int(gelError.ExitCode)
}
