package processors

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type SMA struct {
	period *models.TypedParam[int]
	sma    *indicators.SMA
}

const (
	PeriodName = "period"

	TypeSMA = "SMA"
)

func NewSMA() *SMA {
	periodRange := constraints.NewValRange(1, 200)

	return &SMA{
		period: models.NewParam[int](periodRange),
		sma:    nil,
	}
}

func (sma *SMA) Type() string {
	return TypeSMA
}

func (sma *SMA) Params() models.Params {
	return models.Params{
		PeriodName: sma.period,
	}
}

func (sma *SMA) Initialize() error {
	sma.sma = indicators.NewSMA(sma.period.Value)

	return nil
}

func (sma *SMA) WarmupPeriod() int {
	return sma.sma.Period()
}

func (sma *SMA) Output() float64 {
	return sma.sma.Current()
}

func (sma *SMA) WarmUp(vals []float64) {
	_ = sma.sma.WarmUp(vals)
}

func (sma *SMA) Update(val float64) {
	sma.Update(val)
}
