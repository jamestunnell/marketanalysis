package sources

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/models"
)

type Candlestick struct {
	barValueType *models.TypedParam[string]
	output       float64
}

const TypeCandlestick = "Candlestick"

func NewCandlestick() *Candlestick {
	barValueTypesEnum := constraints.NewValOneOf(BarValueTypes())

	return &Candlestick{
		barValueType: models.NewParam[string](barValueTypesEnum),
		output:       0.0,
	}
}

func (ha *Candlestick) Type() string {
	return TypeCandlestick
}

func (ha *Candlestick) Params() models.Params {
	return models.Params{
		BarValueName: ha.barValueType,
	}
}

func (ha *Candlestick) Initialize() error {
	ha.output = 0.0

	return nil
}

func (ha *Candlestick) WarmupPeriod() int {
	return 1
}

func (ha *Candlestick) Output() float64 {
	return ha.output
}

func (ha *Candlestick) WarmUp(bars models.Bars) error {
	if len(bars) != 1 {
		return commonerrs.NewErrExactLen("warmup bars", len(bars), 1)
	}

	ha.Update(bars[0])

	return nil
}

func (ha *Candlestick) Update(bar *models.Bar) {
	ha.output = BarValue(ha.barValueType.Value, bar.OHLC)
}
