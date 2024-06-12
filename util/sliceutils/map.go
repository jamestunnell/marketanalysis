package sliceutils

func Map[S, T any](ss []S, f func(s S) T) []T {
	ts := make([]T, len(ss))

	for i, s := range ss {
		ts[i] = f(s)
	}

	return ts
}

func MapErr[S, T any](ss []S, f func(s S) (T, error)) ([]T, error) {
	ts := make([]T, len(ss))

	for i, s := range ss {
		t, err := f(s)
		if err != nil {
			return []T{}, err
		}

		ts[i] = t
	}

	return ts, nil
}
