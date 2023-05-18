package sources

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/models"
)

type HeikinAshi struct {
	barValueType *models.TypedParam[string]
	output       float64
	prevOHLC     *models.OHLC
	warm         bool
}

const TypeHeikinAshi = "HeikinAshi"

func NewHeikinAshi() *HeikinAshi {
	barValueTypesEnum := constraints.NewValOneOf(BarValueTypes())

	return &HeikinAshi{
		barValueType: models.NewParam[string](barValueTypesEnum),
		output:       0.0,
		prevOHLC:     nil,
		warm:         false,
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
	ha.warm = false

	return nil
}

func (ha *HeikinAshi) WarmupPeriod() int {
	return 2
}

func (ha *HeikinAshi) Output() float64 {
	return ha.output
}

func (ha *HeikinAshi) Warm() bool {
	return ha.warm
}

func (ha *HeikinAshi) Update(bar *models.Bar) {
	if ha.prevOHLC == nil {
		ha.prevOHLC = bar.OHLC

		return
	}

	haOHLC := bar.HeikinAshi(ha.prevOHLC)

	ha.output = BarValue(ha.barValueType.Value, haOHLC)
	ha.prevOHLC = bar.OHLC
	ha.warm = true
}
