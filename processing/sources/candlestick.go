package sources

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/models"
)

type Candlestick struct {
	barValueType *models.TypedParam[string]
	output       float64
	warm         bool
}

const TypeCandlestick = "Candlestick"

func NewCandlestick() *Candlestick {
	barValueTypesEnum := constraints.NewValOneOf(BarValueTypes())

	return &Candlestick{
		barValueType: models.NewParam[string](barValueTypesEnum),
		output:       0.0,
		warm:         false,
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
	c.warm = false

	return nil
}

func (c *Candlestick) WarmupPeriod() int {
	return 1
}

func (c *Candlestick) Warm() bool {
	return c.warm
}

func (c *Candlestick) Output() float64 {
	return c.output
}

func (c *Candlestick) Update(bar *models.Bar) {
	c.warm = true
	c.output = BarValue(c.barValueType.Value, bar.OHLC)
}
