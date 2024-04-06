package blocks

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMA struct {
	period *models.TypedParam[int]
	ema    *indicators.EMA
	in     *models.TypedInput[float64]
	out    *models.TypedOutput[float64]
}

const (
	DescrEMA = "Exponential moving average"
	TypeEMA = "EMA"
)

func NewEMA() models.Block {
	periodRange := constraints.NewValRange(1, 200)

	return &EMA{
		period: models.NewParam[int](periodRange),
		ema:    nil,
		in:     models.NewTypedInput[float64](),
		out:    models.NewTypedOutput[float64](),
	}
}

func (ma *EMA) GetType() string {
	return TypeEMA
}

func (ma *EMA) GetDescription() string {
	return DescrEMA
}

func (ma *EMA) GetParams() models.Params {
	return models.Params{
		NamePeriod: ma.period,
	}
}

func (ma *EMA) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: ma.in,
	}
}

func (ma *EMA) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOut: ma.out,
	}
}

func (ma *EMA) IsWarm() bool {
	return ma.ema.Warm()
}

func (ma *EMA) Init() error {
	ma.ema = indicators.NewEMA(ma.period.Value)

	return nil
}

func (ma *EMA) Update() {
	if !ma.in.IsSet() {
		return
	}

	ma.ema.Update(ma.in.Get())

	ma.out.Set(ma.ema.Current())
}
