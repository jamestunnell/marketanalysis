package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMA struct {
	len        int
	warm       bool
	current, k float64
}

func NewEMA(len int) models.Indicator {
	return &EMA{
		len:     len,
		warm:    false,
		current: 0.0,
		k:       2.0 / (float64(len) + 1),
	}
}

func (ema *EMA) WarmupPeriod() int {
	return ema.len
}

func (ema *EMA) WarmUp(bars models.Bars) error {
	if len(bars) != ema.len {
		return commonerrs.NewErrExactBarCount("warm up", ema.len, len(bars))
	}

	sum := 0.0
	for _, close := range bars.ClosePrices() {
		sum += close
	}

	ema.current = sum
	ema.warm = true

	return nil
}

func (ema *EMA) Update(bar *models.Bar) float64 {
	if !ema.warm {
		return 0.0
	}

	ema.current = (bar.Close * ema.k) + (ema.current * (1.0 - ema.k))

	return ema.current
}

func (ema *EMA) Current() float64 {
	if !ema.warm {
		return 0.0
	}

	return ema.current
}
