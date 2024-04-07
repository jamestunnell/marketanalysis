package indicators

import (
	"github.com/jamestunnell/marketanalysis/models"
)

type ATR struct {
	current float64
	period  int
	ma      *EMA
	warm    bool
	prev    *models.OHLC
}

func NewATR(period int) *ATR {
	return &ATR{
		current: 0.0,
		period:  period,
		ma:      NewEMA(period),
		warm:    false,
		prev:    nil,
	}
}

func (atr *ATR) Period() int {
	return atr.period
}

func (atr *ATR) Warm() bool {
	return atr.warm
}

func (atr *ATR) Update(cur *models.OHLC) {
	if atr.prev == nil {
		atr.prev = cur

		return
	}

	tr := TrueRange(cur, atr.prev)

	atr.ma.Update(tr)

	if !atr.ma.Warm() {
		return
	}

	atr.current = atr.ma.Current()
}

func (atr *ATR) Current() float64 {
	return atr.current
}
