package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/util/buffer"
)

type Aroon struct {
	period          int
	prev            *buffer.CircularBuffer[float64]
	highAge, lowAge int
	up, down        float64
	warm            bool
}

func NewAroon(period int) *Aroon {
	return &Aroon{
		period:  period,
		prev:    buffer.NewCircularBuffer[float64](period),
		highAge: 0,
		lowAge:  0,
		up:      0.0,
		down:    0.0,
		warm:    false,
	}
}

func (ind *Aroon) WarmupPeriod() int {
	return ind.period
}

func (ind *Aroon) Warm() bool {
	return ind.warm
}

func (ind *Aroon) Update(val float64) {
	ind.prev.Add(val)

	if !ind.prev.Full() {
		return
	}

	min := math.MaxFloat64
	max := -math.MaxFloat64
	highAge := 0
	lowAge := 0
	nMinus1 := ind.period - 1

	ind.prev.EachWithIndex(func(i int, f float64) {
		if f >= max {
			max = f
			highAge = nMinus1 - i
		}

		if f <= min {
			min = f
			lowAge = nMinus1 - i
		}
	})

	n := float64(ind.period)

	ind.warm = true
	ind.up = (n - float64(highAge)) / n
	ind.down = (n - float64(lowAge)) / n
}

func (ind *Aroon) Up() float64 {
	return ind.up
}

func (ind *Aroon) Down() float64 {
	return ind.down
}
