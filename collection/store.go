package collection

type Store interface {
	ItemNames() ([]string, error)
	LoadItem(name string) ([]byte, error)
	StoreItem(name string, data []byte) error
}
