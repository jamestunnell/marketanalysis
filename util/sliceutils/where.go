package sliceutils

func Where[T any](ts []T, f func(t T) bool) []T {
	ts2 := []T{}

	for _, t := range ts {
		if f(t) {
			ts2 = append(ts2, t)
		}
	}

	return ts2
}
