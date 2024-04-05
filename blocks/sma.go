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
	TypeSMA  = "SMA"
	DescrSMA = "Simple moving average"
)

func NewSMA() models.Block {
	periodRange := constraints.NewValRange(1, 200)

	return &SMA{
		period: models.NewParam[int](periodRange),
		sma:    nil,
		in:     models.NewTypedInput[float64](),
		out:    models.NewTypedOutput[float64](),
	}
}

func (sma *SMA) GetType() string {
	return TypeSMA
}

func (sma *SMA) GetDescription() string {
	return DescrSMA
}

func (sma *SMA) GetParams() models.Params {
	return models.Params{
		NamePeriod: sma.period,
	}
}

func (sma *SMA) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: sma.in,
	}
}

func (sma *SMA) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOut: sma.out,
	}
}

func (sma *SMA) Init() error {
	sma.sma = indicators.NewSMA(sma.period.Value)

	return nil
}

func (sma *SMA) IsWarm() bool {
	return sma.sma.Warm()
}

func (sma *SMA) Step() {
	if !sma.in.IsSet() {
		return
	}

	sma.sma.Update(sma.in.Get())

	sma.out.Set(sma.sma.Current())
}
