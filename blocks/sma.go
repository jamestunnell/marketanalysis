package blocks

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type SMA struct {
	period *models.TypedParam[int]
	sma    *indicators.SMA
	in     *models.TypedInput[float64]
	out    *models.TypedOutput[float64]
}

const (
	DescrSMA = "Simple Moving Average"
	TypeSMA  = "SMA"
)

func NewSMA() models.Block {
	periodRange := constraints.NewValRange(1, 200)

	return &SMA{
		period: models.NewParam[int](1, periodRange),
		sma:    nil,
		in:     models.NewTypedInput[float64](),
		out:    models.NewTypedOutput[float64](),
	}
}

func (blk *SMA) GetType() string {
	return TypeSMA
}

func (blk *SMA) GetDescription() string {
	return DescrSMA
}

func (blk *SMA) GetParams() models.Params {
	return models.Params{
		NamePeriod: blk.period,
	}
}

func (blk *SMA) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: blk.in,
	}
}

func (blk *SMA) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOut: blk.out,
	}
}

func (blk *SMA) IsWarm() bool {
	return blk.sma.Warm()
}

func (blk *SMA) Init() error {
	blk.sma = indicators.NewSMA(blk.period.Value)

	return nil
}

func (blk *SMA) Update() {
	if !blk.in.IsSet() {
		return
	}

	blk.sma.Update(blk.in.Get())

	blk.out.Set(blk.sma.Current())
}
