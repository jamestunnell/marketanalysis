package logic

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Xor struct {
	in1 *blocks.TypedInput[bool]
	in2 *blocks.TypedInput[bool]
	out *blocks.TypedOutput[bool]
}

const (
	TypeXor  = "Xor"
	DescrXor = "XOR two bool signals."
)

func NewXor() blocks.Block {
	return &Xor{
		in1: blocks.NewTypedInput[bool](),
		in2: blocks.NewTypedInput[bool](),
		out: blocks.NewTypedOutput[bool](),
	}
}

func (blk *Xor) GetType() string {
	return TypeXor
}

func (blk *Xor) GetDescription() string {
	return DescrXor
}

func (blk *Xor) GetParams() models.Params {
	return models.Params{}
}

func (blk *Xor) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Xor) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Xor) GetWarmupPeriod() int {
	return 0
}

func (blk *Xor) IsWarm() bool {
	return true
}

func (blk *Xor) Init() error {
	return nil
}

func (blk *Xor) Update(_ *models.Bar, isLast bool) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() != blk.in2.GetValue())
}
