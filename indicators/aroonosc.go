package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/util/buffer"
)

type AroonOsc struct {
	period          int
	prev            *buffer.CircularBuffer[float64]
	highAge, lowAge int
	up, down, diff  float64
}

func NewAroonOsc(period int) *AroonOsc {
	return &AroonOsc{
		period:  period,
		prev:    buffer.NewCircularBuffer[float64](period),
		highAge: 0,
		lowAge:  0,
		up:      0.0,
		down:    0.0,
		diff:    0.0,
	}
}

func (osc *AroonOsc) WarmupPeriod() int {
	return osc.period
}

func (osc *AroonOsc) WarmUp(vals []float64) error {
	osc.prev.AddN(vals[:osc.WarmupPeriod()-1]...)

	for i := osc.WarmupPeriod() - 1; i < len(vals); i++ {
		osc.Update(vals[i])
	}

	return nil
}

func (osc *AroonOsc) Update(val float64) {
	osc.prev.Add(val)

	min := math.MaxFloat64
	max := -math.MaxFloat64
	highAge := 0
	lowAge := 0
	nMinus1 := osc.period - 1

	osc.prev.EachWithIndex(func(i int, f float64) {
		if f >= max {
			max = f
			highAge = nMinus1 - i
		}

		if f <= min {
			min = f
			lowAge = nMinus1 - i
		}
	})

	n := float64(osc.period)

	osc.up = (n - float64(highAge)) / n
	osc.down = (n - float64(lowAge)) / n
	osc.diff = osc.up - osc.down
}

func (osc *AroonOsc) Up() float64 {
	return osc.up
}

func (osc *AroonOsc) Down() float64 {
	return osc.down
}

func (osc *AroonOsc) Diff() float64 {
	return osc.diff
}
