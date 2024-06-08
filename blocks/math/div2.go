package math

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Div2 struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
}

const (
	TypeDiv2  = "Div2"
	DescrDiv2 = "Divide one input signal by another."
)

func NewDiv2() blocks.Block {
	return &Div2{
		in1: blocks.NewTypedInput[float64](),
		in2: blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *Div2) GetType() string {
	return TypeDiv2
}

func (blk *Div2) GetDescription() string {
	return DescrDiv2
}

func (blk *Div2) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *Div2) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Div2) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Div2) GetWarmupPeriod() int {
	return 0
}

func (blk *Div2) IsWarm() bool {
	return true
}

func (blk *Div2) Init() error {
	return nil
}

func (blk *Div2) Update(_ *models.Bar, isLast bool) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	// In case the divisor is 0, Go will produce +Inf, -Inf, or NaN.
	blk.out.SetValue(blk.in1.GetValue() / blk.in2.GetValue())
}
