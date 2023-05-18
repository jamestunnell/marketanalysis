package processors

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMA struct {
	period *models.TypedParam[int]
	ema    *indicators.EMA
}

const TypeEMA = "EMA"

func NewEMA() *EMA {
	periodRange := constraints.NewValRange(1, 200)

	return &EMA{
		period: models.NewParam[int](periodRange),
		ema:    nil,
	}
}

func (ema *EMA) Type() string {
	return TypeEMA
}

func (ema *EMA) Params() models.Params {
	return models.Params{
		PeriodName: ema.period,
	}
}

func (ema *EMA) Initialize() error {
	ema.ema = indicators.NewEMA(ema.period.Value)

	return nil
}

func (ema *EMA) WarmupPeriod() int {
	return ema.ema.Period()
}

func (ema *EMA) Output() float64 {
	return ema.ema.Current()
}

func (ema *EMA) Warm() bool {
	return ema.ema.Warm()
}

func (ema *EMA) Update(val float64) {
	ema.ema.Update(val)
}
