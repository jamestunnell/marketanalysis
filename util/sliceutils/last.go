package sliceutils

func Last[T any](ts []T) T {
	return ts[len(ts)-1]
}

func LastN[T any](ts []T, n int) []T {
	return ts[len(ts)-n:]
}
