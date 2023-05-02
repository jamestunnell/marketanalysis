package sources

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type TrueRange struct {
	prevOHLC *models.OHLC
	output   float64
}

const TypeTrueRange = "TrueRange"

func NewTrueRange() *TrueRange {
	return &TrueRange{
		output:   0.0,
		prevOHLC: nil,
	}
}

func (ha *TrueRange) Type() string {
	return TypeTrueRange
}

func (ha *TrueRange) Params() models.Params {
	return models.Params{}
}

func (ha *TrueRange) Initialize() error {
	ha.output = 0.0
	ha.prevOHLC = nil

	return nil
}

func (ha *TrueRange) WarmupPeriod() int {
	return 2
}

func (ha *TrueRange) Output() float64 {
	return ha.output
}

func (ha *TrueRange) WarmUp(bars models.Bars) error {
	if len(bars) != 2 {
		return commonerrs.NewErrExactLen("warmup bars", len(bars), 2)
	}

	ha.prevOHLC = bars[0].OHLC

	ha.Update(bars[1])

	return nil
}

func (ha *TrueRange) Update(bar *models.Bar) {
	ha.output = indicators.TrueRange(bar.OHLC, ha.prevOHLC)
	ha.prevOHLC = bar.OHLC
}
