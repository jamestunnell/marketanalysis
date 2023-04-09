package buffer

import "golang.org/x/exp/slices"

type FullCircularBuffer struct {
	elems          []float64
	length, oldest int
}

func NewFullCircularBuffer(startElems []float64) *FullCircularBuffer {
	return &FullCircularBuffer{
		elems:  slices.Clone(startElems),
		length: len(startElems),
		oldest: 0,
	}
}

func (cb *FullCircularBuffer) Add(x float64) {
	cb.elems[cb.oldest] = x

	cb.oldest++

	if cb.oldest == cb.length {
		cb.oldest = 0
	}
}

// func (cb *FullCircularBuffer) Front(n int) []float64 {

// }

// func (cb *FullCircularBuffer) Back(n int) []float64 {

// }

func (cb *FullCircularBuffer) Length() int {
	return cb.length
}

func (cb *FullCircularBuffer) Sum() float64 {
	sum := 0.0

	// order doesn't matter for sum
	for _, elem := range cb.elems {
		sum += elem
	}

	return sum
}
