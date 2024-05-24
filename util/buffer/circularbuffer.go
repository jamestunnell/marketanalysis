package buffer

type CircularBuffer[T any] struct {
	elems           []T
	cap, head, tail int
	empty, full     bool
}

func NewCircularBuffer[T any](cap int) *CircularBuffer[T] {
	return &CircularBuffer[T]{
		elems: make([]T, cap),
		cap:   cap,
		head:  0,
		tail:  0,
		empty: true,
		full:  false,
	}
}

func (cb *CircularBuffer[T]) Capacity() int {
	return cb.cap
}

func (cb *CircularBuffer[T]) Len() int {
	switch {
	case cb.full:
		return cb.cap
	case cb.empty:
		return 0
	case cb.head > cb.tail:
		return cb.head - cb.tail
	}

	return cb.cap - (cb.tail - cb.head)
}

func (cb *CircularBuffer[T]) Empty() bool {
	return cb.empty
}

func (cb *CircularBuffer[T]) Full() bool {
	return cb.full
}

func (cb *CircularBuffer[T]) AddN(ts ...T) {
	for _, t := range ts {
		cb.Add(t)
	}
}

func (cb *CircularBuffer[T]) Add(t T) {
	if cb.full {
		cb.elems[cb.head] = t

		cb.incrHead()
		cb.incrTail()

		return
	}

	cb.elems[cb.head] = t

	cb.empty = false
	cb.incrHead()

	if cb.head == cb.tail {
		cb.full = true
	}
}

func (cb *CircularBuffer[T]) Newest() (T, bool) {
	var t T

	if cb.empty {
		return t, false
	}

	return cb.elems[(cb.tail+cb.Len()-1)%cb.cap], true
}

func (cb *CircularBuffer[T]) Array() []T {
	ts := []T{}

	cb.Each(func(t T) {
		ts = append(ts, t)
	})

	return ts
}

func (cb *CircularBuffer[T]) Each(f func(T)) {
	for i := cb.tail; i < (cb.tail + cb.Len()); i++ {
		f(cb.elems[i%cb.cap])
	}
}

func (cb *CircularBuffer[T]) EachWithIndex(f func(int, T)) {
	apparentIndex := 0
	for i := cb.tail; i < (cb.tail + cb.Len()); i++ {
		f(apparentIndex, cb.elems[i%cb.cap])

		apparentIndex++
	}
}

func (cb *CircularBuffer[T]) incrHead() {
	cb.head++

	if cb.head >= cb.cap {
		cb.head = 0
	}
}

func (cb *CircularBuffer[T]) incrTail() {
	cb.tail++

	if cb.tail >= cb.cap {
		cb.tail = 0
	}
}
