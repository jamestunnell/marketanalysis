package collection

type Store interface {
	MakeSubstore(name string) error
	SubstoreNames() []string
	ItemNames() []string
	Substore(name string) (Store, error)
	Item(name string) (Item, error)
}
