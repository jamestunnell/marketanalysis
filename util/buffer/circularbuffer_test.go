package buffer_test

import (
	"testing"

	"github.com/jamestunnell/marketanalysis/util/buffer"
	"github.com/stretchr/testify/assert"
)

func TestCircularBuffer_Add(t *testing.T) {
	cb := buffer.NewCircularBuffer[float64](5)

	assert.Empty(t, cb.Array())

	cb.Add(1.0)

	assert.Equal(t, []float64{1.0}, cb.Array())

	cb.Add(2.0)

	assert.Equal(t, []float64{1.0, 2.0}, cb.Array())

	cb.AddN(3.0, 4.0, 5.0, 6.0)

	assert.Equal(t, []float64{2.0, 3.0, 4.0, 5.0, 6.0}, cb.Array())
}

func TestCircularBuffer_EachWithIndex(t *testing.T) {
	cb := buffer.NewCircularBuffer[float64](5)

	cb.AddN(1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0)

	indices := []int{}
	values := []float64{}

	cb.EachWithIndex(func(i int, f float64) {
		indices = append(indices, i)
		values = append(values, f)
	})

	assert.Equal(t, []int{0, 1, 2, 3, 4}, indices)
	assert.Equal(t, []float64{3.0, 4.0, 5.0, 6.0, 7.0}, values)
}
