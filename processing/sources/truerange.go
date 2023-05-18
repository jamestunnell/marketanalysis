package sources

import (
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type TrueRange struct {
	prevOHLC *models.OHLC
	output   float64
	warm     bool
}

const TypeTrueRange = "TrueRange"

func NewTrueRange() *TrueRange {
	return &TrueRange{
		output:   0.0,
		prevOHLC: nil,
		warm:     false,
	}
}

func (tr *TrueRange) Type() string {
	return TypeTrueRange
}

func (tr *TrueRange) Params() models.Params {
	return models.Params{}
}

func (tr *TrueRange) Initialize() error {
	tr.output = 0.0
	tr.prevOHLC = nil
	tr.warm = false

	return nil
}

func (tr *TrueRange) WarmupPeriod() int {
	return 2
}

func (tr *TrueRange) Output() float64 {
	return tr.output
}

func (tr *TrueRange) Warm() bool {
	return tr.warm
}

func (tr *TrueRange) Update(bar *models.Bar) {
	if tr.prevOHLC == nil {
		tr.prevOHLC = bar.OHLC

		return
	}

	tr.output = indicators.TrueRange(bar.OHLC, tr.prevOHLC)
	tr.prevOHLC = bar.OHLC
	tr.warm = true
}
