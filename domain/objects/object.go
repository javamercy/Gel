package objects

type IObject interface {
	Type() string
	Size() int64
	Content() []byte
	Header() []byte
	Sha() string
}
