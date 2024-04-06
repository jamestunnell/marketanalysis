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
	DescrSMA = "Simple moving average"
	TypeSMA  = "SMA"
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

func (ma *SMA) GetType() string {
	return TypeSMA
}

func (ma *SMA) GetDescription() string {
	return DescrSMA
}

func (ma *SMA) GetParams() models.Params {
	return models.Params{
		NamePeriod: ma.period,
	}
}

func (ma *SMA) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: ma.in,
	}
}

func (ma *SMA) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOut: ma.out,
	}
}

func (ma *SMA) IsWarm() bool {
	return ma.sma.Warm()
}

func (ma *SMA) Init() error {
	ma.sma = indicators.NewSMA(ma.period.Value)

	return nil
}

func (ma *SMA) Update() {
	if !ma.in.IsSet() {
		return
	}

	ma.sma.Update(ma.in.Get())

	ma.out.Set(ma.sma.Current())
}
