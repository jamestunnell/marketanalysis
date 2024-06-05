package math

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type MulK struct {
	in  *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
	k   *blocks.FloatParam
}

const (
	TypeMulK  = "MulK"
	DescrMulK = "Multiply an input signal by a constant."
)

func NewMulK() blocks.Block {
	return &MulK{
		in:  blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
		k:   blocks.NewFloatParam(1.0, blocks.NewNone[float64]()),
	}
}

func (blk *MulK) GetType() string {
	return TypeMulK
}

func (blk *MulK) GetDescription() string {
	return DescrMulK
}

func (blk *MulK) GetParams() blocks.Params {
	return blocks.Params{
		blocks.NameK: blk.k,
	}
}

func (blk *MulK) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *MulK) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *MulK) GetWarmupPeriod() int {
	return 0
}

func (blk *MulK) IsWarm() bool {
	return true
}

func (blk *MulK) Init() error {
	return nil
}

func (blk *MulK) Update(_ *models.Bar) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in.GetValue() * blk.k.CurrentVal)
}
