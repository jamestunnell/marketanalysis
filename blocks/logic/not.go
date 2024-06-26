package logic

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Not struct {
	in  *blocks.TypedInput[bool]
	out *blocks.TypedOutput[bool]
}

const (
	TypeNot  = "Not"
	DescrNot = "NOT a bool signal."
)

func NewNot() blocks.Block {
	return &Not{
		in:  blocks.NewTypedInput[bool](),
		out: blocks.NewTypedOutput[bool](),
	}
}

func (blk *Not) GetType() string {
	return TypeNot
}

func (blk *Not) GetDescription() string {
	return DescrNot
}

func (blk *Not) GetParams() models.Params {
	return models.Params{}
}

func (blk *Not) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *Not) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Not) GetWarmupPeriod() int {
	return 0
}

func (blk *Not) IsWarm() bool {
	return true
}

func (blk *Not) Init() error {
	return nil
}

func (blk *Not) Update(_ *models.Bar, isLast bool) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.out.SetValue(!blk.in.GetValue())
}
