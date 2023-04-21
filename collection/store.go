package collection

type Store interface {
	MakeSubstore(name string) (Store, error)
	SubstoreNames() []string
	Substore(name string) (Store, error)

	ItemNames() []string
	LoadItem(name string) ([]byte, error)
	StoreItem(name string, data []byte) error
}
