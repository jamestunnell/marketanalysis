package collection

type Collection struct {
	Store Store
}

func Load(store Store) *Collection {
	return &Collection{
		Store: store,
	}
}
