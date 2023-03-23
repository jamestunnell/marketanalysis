package collection

type Collection struct {
	Info *Info
}

func Load(store Store) *Collection {
	return &Collection{
		Info: store.GetInfo(),
	}
}
