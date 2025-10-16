package repositories

type IRepository interface {
	MakeDir(path string) error
	WriteFile(path string, data []byte) error
	Exists(path string) bool
	ReadFile(path string) ([]byte, error)
}
