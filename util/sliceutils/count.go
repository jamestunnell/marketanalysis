package sliceutils

func Count[T any](ts []T, f func(t T) bool) int {
	count := 0

	for _, t := range ts {
		if f(t) {
			count++
		}
	}

	return count
}
