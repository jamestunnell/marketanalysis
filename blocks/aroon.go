package blocks

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Aroon struct {
	period *models.TypedParam[int]
	aroon  *indicators.Aroon
	in     *models.TypedInput[float64]
	up     *models.TypedOutput[float64]
	dn     *models.TypedOutput[float64]
}

const (
	DescrAroon = "Aroon indicator identifies trend changes in the price of an asset, as well as the strength of that trend."
	TypeAroon  = "Aroon"
)

func NewAroon() models.Block {
	periodRange := constraints.NewValRange(1, 200)

	return &Aroon{
		period: models.NewParam[int](periodRange),
		aroon:  nil,
		in:     models.NewTypedInput[float64](),
		up:     models.NewTypedOutput[float64](),
		dn:     models.NewTypedOutput[float64](),
	}
}

func (blk *Aroon) GetType() string {
	return TypeAroon
}

func (blk *Aroon) GetDescription() string {
	return DescrAroon
}

func (blk *Aroon) GetParams() models.Params {
	return models.Params{
		NamePeriod: blk.period,
	}
}

func (blk *Aroon) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: blk.in,
	}
}

func (blk *Aroon) GetOutputs() models.Outputs {
	return models.Outputs{
		NameUp:   blk.up,
		NameDown: blk.dn,
	}
}

func (blk *Aroon) IsWarm() bool {
	return blk.aroon.Warm()
}

func (blk *Aroon) Init() error {
	blk.aroon = indicators.NewAroon(blk.period.Value)

	return nil
}

func (blk *Aroon) Update() {
	if !blk.in.IsSet() {
		return
	}

	blk.aroon.Update(blk.in.Get())

	if !blk.aroon.Warm() {
		return
	}

	blk.up.Set(blk.aroon.Up())
	blk.dn.Set(blk.aroon.Down())
}
