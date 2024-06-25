package math

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Mul2 struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
}

const (
	TypeMul2  = "Mul2"
	DescrMul2 = "Multiply two input signals."

	NameIn1 = blocks.NameIn + "1"
	NameIn2 = blocks.NameIn + "2"
)

func NewMul2() blocks.Block {
	return &Mul2{
		in1: blocks.NewTypedInput[float64](),
		in2: blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *Mul2) GetType() string {
	return TypeMul2
}

func (blk *Mul2) GetDescription() string {
	return DescrMul2
}

func (blk *Mul2) GetParams() models.Params {
	return models.Params{}
}

func (blk *Mul2) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Mul2) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Mul2) GetWarmupPeriod() int {
	return 0
}

func (blk *Mul2) IsWarm() bool {
	return true
}

func (blk *Mul2) Init() error {
	return nil
}

func (blk *Mul2) Update(_ *models.Bar, isLast bool) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() * blk.in2.GetValue())
}
