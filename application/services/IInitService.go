package services

type IInitService interface {
	Init(path string) error
}
