package processing

type ElementRegistry[T Element] struct {
	elems map[string]func() T
}

func NewElementRegistry[T Element]() *ElementRegistry[T] {
	return &ElementRegistry[T]{
		elems: map[string]func() T{},
	}
}

func (r *ElementRegistry[T]) Add(typ string, newElem func() T) {
	r.elems[typ] = newElem
}

func (r *ElementRegistry[T]) Get(typ string) (func() T, bool) {
	f, found := r.elems[typ]
	if !found {
		return nil, false
	}

	return f, true
}
