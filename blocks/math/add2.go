package math

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Add2 struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
}

const (
	TypeAdd2  = "Add2"
	DescrAdd2 = "Add two input signals."
)

func NewAdd2() blocks.Block {
	return &Add2{
		in1: blocks.NewTypedInput[float64](),
		in2: blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *Add2) GetType() string {
	return TypeAdd2
}

func (blk *Add2) GetDescription() string {
	return DescrAdd2
}

func (blk *Add2) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *Add2) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Add2) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Add2) GetWarmupPeriod() int {
	return 0
}

func (blk *Add2) IsWarm() bool {
	return true
}

func (blk *Add2) Init() error {
	return nil
}

func (blk *Add2) Update(_ *models.Bar, isLast bool) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() + blk.in2.GetValue())
}
