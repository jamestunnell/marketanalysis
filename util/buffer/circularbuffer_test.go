package buffer_test

import (
	"testing"

	"github.com/jamestunnell/marketanalysis/util/buffer"
	"github.com/stretchr/testify/assert"
)

func TestCircularBuffer(t *testing.T) {
	cb := buffer.NewCircularBuffer[float64](5)

	assert.Empty(t, cb.Array())

	cb.Add(1.0)

	assert.Equal(t, []float64{1.0}, cb.Array())

	cb.Add(2.0)

	assert.Equal(t, []float64{1.0, 2.0}, cb.Array())

	cb.AddN(3.0, 4.0, 5.0, 6.0)

	assert.Equal(t, []float64{2.0, 3.0, 4.0, 5.0, 6.0}, cb.Array())
}
