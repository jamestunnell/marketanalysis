package classification_test

import (
	"testing"
	"time"

	"github.com/jamestunnell/marketanalysis/classification"
	"github.com/jamestunnell/marketanalysis/models/bar"
	"github.com/stretchr/testify/assert"
)

func TestMeasureBar(t *testing.T) {
	b := bar.New(time.Now(), 20.0, 26.2, 18.5, 22.0, 12, 12, 20.0)
	m := classification.MeasureBar(b)

	assert.InDelta(t, 2.0, m.Body, 1e-10)
	assert.InDelta(t, 4.2, m.TopWick, 1e-10)
	assert.InDelta(t, 1.5, m.BottomWick, 1e-10)
	assert.True(t, m.Bullish)

	// swap open/close to make it bearish
	b.Open, b.Close = b.Close, b.Open

	m = classification.MeasureBar(b)

	assert.InDelta(t, 2.0, m.Body, 1e-10)
	assert.InDelta(t, 4.2, m.TopWick, 1e-10)
	assert.InDelta(t, 1.5, m.BottomWick, 1e-10)
	assert.False(t, m.Bullish)
}
