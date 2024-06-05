package math

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type AddK struct {
	in  *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
	k   *blocks.FloatParam
}

const (
	TypeAddK  = "AddK"
	DescrAddK = "Add an input signal and a constant."
)

func NewAddK() blocks.Block {
	return &AddK{
		in:  blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
		k:   blocks.NewFloatParam(0.0, blocks.NewNone[float64]()),
	}
}

func (blk *AddK) GetType() string {
	return TypeAddK
}

func (blk *AddK) GetDescription() string {
	return DescrAddK
}

func (blk *AddK) GetParams() blocks.Params {
	return blocks.Params{
		blocks.NameK: blk.k,
	}
}

func (blk *AddK) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *AddK) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *AddK) GetWarmupPeriod() int {
	return 0
}

func (blk *AddK) IsWarm() bool {
	return true
}

func (blk *AddK) Init() error {
	return nil
}

func (blk *AddK) Update(_ *models.Bar) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in.GetValue() + blk.k.CurrentVal)
}
