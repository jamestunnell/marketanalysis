package collection

type Item interface {
	Name() string
	Load() ([]byte, error)
	Store([]byte) error
}
