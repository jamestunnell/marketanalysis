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
	DescrEMA = "Exponential Moving Average"
	TypeEMA  = "EMA"
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

func (blk *EMA) GetType() string {
	return TypeEMA
}

func (blk *EMA) GetDescription() string {
	return DescrEMA
}

func (blk *EMA) GetParams() models.Params {
	return models.Params{
		NamePeriod: blk.period,
	}
}

func (blk *EMA) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: blk.in,
	}
}

func (blk *EMA) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOut: blk.out,
	}
}

func (blk *EMA) IsWarm() bool {
	return blk.ema.Warm()
}

func (blk *EMA) Init() error {
	blk.ema = indicators.NewEMA(blk.period.Value)

	return nil
}

func (blk *EMA) Update() {
	if !blk.in.IsSet() {
		return
	}

	blk.ema.Update(blk.in.Get())

	blk.out.Set(blk.ema.Current())
}
