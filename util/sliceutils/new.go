package sliceutils

func New[T any](len int, f func(idx int) T) []T {
	ts := make([]T, len)

	for i := 0; i < len; i++ {
		ts[i] = f(i)
	}

	return ts
}
