package logic

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type And struct {
	in1 *blocks.TypedInput[bool]
	in2 *blocks.TypedInput[bool]
	out *blocks.TypedOutput[bool]
}

const (
	TypeAnd  = "And"
	DescrAnd = "AND two bool signals."

	NameIn1 = "in1"
	NameIn2 = "in2"
)

func NewAnd() blocks.Block {
	return &And{
		in1: blocks.NewTypedInput[bool](),
		in2: blocks.NewTypedInput[bool](),
		out: blocks.NewTypedOutput[bool](),
	}
}

func (blk *And) GetType() string {
	return TypeAnd
}

func (blk *And) GetDescription() string {
	return DescrAnd
}

func (blk *And) GetParams() models.Params {
	return models.Params{}
}

func (blk *And) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *And) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *And) GetWarmupPeriod() int {
	return 0
}

func (blk *And) IsWarm() bool {
	return true
}

func (blk *And) Init() error {
	return nil
}

func (blk *And) Update(_ *models.Bar, isLast bool) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() && blk.in2.GetValue())
}
