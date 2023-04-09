package prediction_test

import (
	"testing"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/prediction"
	"github.com/stretchr/testify/assert"
)

func TestMeasureBar(t *testing.T) {
	b := models.NewBar(time.Now(), 20.0, 26.2, 18.5, 22.0)
	m := prediction.NewBarMeasure(b, 1.0)

	assert.InDelta(t, 2.0, m.Body, 1e-10)
	assert.InDelta(t, 4.2, m.Top, 1e-10)
	assert.InDelta(t, 1.5, m.Bottom, 1e-10)

	// swap open/close to make it bearish
	b.Open, b.Close = b.Close, b.Open

	m = prediction.NewBarMeasure(b, 1.0)

	assert.InDelta(t, -2.0, m.Body, 1e-10)
	assert.InDelta(t, 4.2, m.Top, 1e-10)
	assert.InDelta(t, 1.5, m.Bottom, 1e-10)
}
