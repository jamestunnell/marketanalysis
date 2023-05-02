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

func (c *Candlestick) Type() string {
	return TypeCandlestick
}

func (c *Candlestick) Params() models.Params {
	return models.Params{
		BarValueName: c.barValueType,
	}
}

func (c *Candlestick) Initialize() error {
	c.output = 0.0

	return nil
}

func (c *Candlestick) WarmupPeriod() int {
	return 1
}

func (c *Candlestick) Output() float64 {
	return c.output
}

func (c *Candlestick) WarmUp(bars models.Bars) error {
	if len(bars) != 1 {
		return commonerrs.NewErrExactLen("warmup bars", len(bars), 1)
	}

	c.Update(bars[0])

	return nil
}

func (c *Candlestick) Update(bar *models.Bar) {
	c.output = BarValue(c.barValueType.Value, bar.OHLC)
}
