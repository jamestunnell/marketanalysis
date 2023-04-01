package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models/bar"
)

type ATR struct {
	length  int
	prevBar *bar.Bar
	current float64
}

func NewATR(length int) *ATR {
	return &ATR{length: length}
}

func (atr *ATR) Length() int {
	return atr.length
}

func (atr *ATR) WarmupPeriod() int {
	return atr.length + 1
}

func (atr *ATR) WarmUp(bars []*bar.Bar) error {
	wp := atr.WarmupPeriod()
	if len(bars) != wp {
		return commonerrs.NewErrExactBarCount("warmup", wp, len(bars))
	}

	sum := 0.0

	atr.prevBar = bars[0]

	for i := 1; i < len(bars); i++ {
		sum += TrueRange(bars[i], bars[i-1])
	}

	atr.prevBar = bars[len(bars)-1]
	atr.current = sum / float64(atr.length)

	return nil
}

func (atr *ATR) Current() float64 {
	return atr.current
}

func (atr *ATR) Update(bar *bar.Bar) float64 {
	tr := TrueRange(bar, atr.prevBar)
	n := float64(atr.length)

	atr.current = (atr.current*(n-1) + tr) / n
	atr.prevBar = bar

	return atr.current
}
