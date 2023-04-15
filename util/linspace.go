package util

import (
	"math"

	"gonum.org/v1/gonum/floats"
)

func LinSpaceInts(start, stop, n int) []int {
	ints := make([]int, n)

	for i, f := range floats.Span(make([]float64, n), float64(start), float64(stop)) {
		ints[i] = int(math.Round(f))
	}

	return ints
}
