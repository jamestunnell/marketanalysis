package sliceutils

func Each[T any](ts []T, f func(t T)) {
	for _, t := range ts {
		f(t)
	}
}

func EachWithIndex[T any](ts []T, f func(t T, idx int)) {
	for i, t := range ts {
		f(t, i)
	}
}
