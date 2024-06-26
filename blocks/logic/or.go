package logic

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Or struct {
	in1 *blocks.TypedInput[bool]
	in2 *blocks.TypedInput[bool]
	out *blocks.TypedOutput[bool]
}

const (
	TypeOr  = "Or"
	DescrOr = "OR two bool signals."
)

func NewOr() blocks.Block {
	return &Or{
		in1: blocks.NewTypedInput[bool](),
		in2: blocks.NewTypedInput[bool](),
		out: blocks.NewTypedOutput[bool](),
	}
}

func (blk *Or) GetType() string {
	return TypeOr
}

func (blk *Or) GetDescription() string {
	return DescrOr
}

func (blk *Or) GetParams() models.Params {
	return models.Params{}
}

func (blk *Or) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Or) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Or) GetWarmupPeriod() int {
	return 0
}

func (blk *Or) IsWarm() bool {
	return true
}

func (blk *Or) Init() error {
	return nil
}

func (blk *Or) Update(_ *models.Bar, isLast bool) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() || blk.in2.GetValue())
}
