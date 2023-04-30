package sources

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/models"
)

type HeikinAshi struct {
	barValueType *models.TypedParam[string]
	output       float64
	prevOHLC     *models.OHLC
}

const TypeHeikinAshi = "HeikinAshi"

func NewHeikinAshi() *HeikinAshi {
	barValueTypesEnum := constraints.NewValOneOf(BarValueTypes())

	return &HeikinAshi{
		barValueType: models.NewParam[string](barValueTypesEnum),
		output:       0.0,
		prevOHLC:     nil,
	}
}

func (ha *HeikinAshi) Type() string {
	return TypeHeikinAshi
}

func (ha *HeikinAshi) Params() models.Params {
	return models.Params{
		BarValueName: ha.barValueType,
	}
}

func (ha *HeikinAshi) Initialize() error {
	ha.output = 0.0
	ha.prevOHLC = nil

	return nil
}

func (ha *HeikinAshi) WarmupPeriod() int {
	return 2
}

func (ha *HeikinAshi) Output() float64 {
	return ha.output
}

func (ha *HeikinAshi) WarmUp(bars models.Bars) {
	ha.prevOHLC = bars[0].OHLC

	ha.Update(bars[1])
}

func (ha *HeikinAshi) Update(bar *models.Bar) {
	ohlc := bar.HeikinAshi(bar.OHLC)

	ha.output = BarValue(ha.barValueType.Value, ohlc)
	ha.prevOHLC = ohlc
}
