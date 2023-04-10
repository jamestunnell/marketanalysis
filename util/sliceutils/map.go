package sliceutils

func Map[S, T any](ss []S, f func(s S) T) []T {
	ts := make([]T, len(ss))

	for i, s := range ss {
		ts[i] = f(s)
	}

	return ts
}
