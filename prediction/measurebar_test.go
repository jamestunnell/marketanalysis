package prediction_test

import (
	"testing"
	"time"

	"github.com/jamestunnell/marketanalysis/models/bar"
	"github.com/jamestunnell/marketanalysis/prediction"
	"github.com/stretchr/testify/assert"
)

func TestMeasureBar(t *testing.T) {
	b := bar.New(time.Now(), 20.0, 26.2, 18.5, 22.0, 12, 12, 20.0)
	body, top, bottom := prediction.MeasureBar(b)

	assert.InDelta(t, 2.0, body, 1e-10)
	assert.InDelta(t, 4.2, top, 1e-10)
	assert.InDelta(t, 1.5, bottom, 1e-10)

	// swap open/close to make it bearish
	b.Open, b.Close = b.Close, b.Open

	body, top, bottom = prediction.MeasureBar(b)

	assert.InDelta(t, -2.0, body, 1e-10)
	assert.InDelta(t, 4.2, top, 1e-10)
	assert.InDelta(t, 1.5, bottom, 1e-10)
}
