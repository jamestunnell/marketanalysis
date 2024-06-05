package math

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Sub2 struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
}

const (
	TypeSub2  = "Sub2"
	DescrSub2 = "Subtract one input signal from another."
)

func NewSub2() blocks.Block {
	return &Sub2{
		in1: blocks.NewTypedInput[float64](),
		in2: blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *Sub2) GetType() string {
	return TypeSub2
}

func (blk *Sub2) GetDescription() string {
	return DescrSub2
}

func (blk *Sub2) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *Sub2) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Sub2) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Sub2) GetWarmupPeriod() int {
	return 0
}

func (blk *Sub2) IsWarm() bool {
	return true
}

func (blk *Sub2) Init() error {
	return nil
}

func (blk *Sub2) Update(_ *models.Bar) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() - blk.in2.GetValue())
}
