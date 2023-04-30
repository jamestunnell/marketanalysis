package processors

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type SMA struct {
	period *models.TypedParam[int]
	ema    *indicators.SMA
}

const TypeSMA = "SMA"

func NewSMA() *SMA {
	periodRange := constraints.NewValRange(1, 200)

	return &SMA{
		period: models.NewParam[int](periodRange),
		ema:    nil,
	}
}

func (ema *SMA) Type() string {
	return TypeSMA
}

func (ema *SMA) Params() models.Params {
	return models.Params{}
}

func (ema *SMA) Initialize() error {
	ema.ema = indicators.NewSMA(ema.period.Value)

	return nil
}

func (ema *SMA) WarmupPeriod() int {
	return ema.ema.Period()
}

func (ema *SMA) Output() float64 {
	return ema.ema.Current()
}

func (ema *SMA) WarmUp(vals []float64) {
	_ = ema.ema.WarmUp(vals)
}

func (ema *SMA) Update(val float64) {
	ema.Update(val)
}
