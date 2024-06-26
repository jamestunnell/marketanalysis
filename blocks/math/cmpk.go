package math

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type CmpK struct {
	in      *blocks.TypedInput[float64]
	less    *blocks.TypedOutput[bool]
	greater *blocks.TypedOutput[bool]
	k       *models.FloatParam
}

const (
	TypeCmpK  = "CmpK"
	DescrCmpK = "Compare an input signal with a constant."
)

func NewCmpK() blocks.Block {
	return &CmpK{
		in:      blocks.NewTypedInput[float64](),
		less:    blocks.NewTypedOutput[bool](),
		greater: blocks.NewTypedOutput[bool](),
		k:       models.NewFloatParam(0.0, models.NewUnconstrained[float64]()),
	}
}

func (blk *CmpK) GetType() string {
	return TypeCmpK
}

func (blk *CmpK) GetDescription() string {
	return DescrCmpK
}

func (blk *CmpK) GetParams() models.Params {
	return models.Params{
		blocks.NameK: blk.k,
	}
}

func (blk *CmpK) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *CmpK) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		"less":    blk.less,
		"greater": blk.greater,
	}
}

func (blk *CmpK) GetWarmupPeriod() int {
	return 0
}

func (blk *CmpK) IsWarm() bool {
	return true
}

func (blk *CmpK) Init() error {
	return nil
}

func (blk *CmpK) Update(_ *models.Bar, isLast bool) {
	if !blk.in.IsValueSet() {
		return
	}

	in := blk.in.GetValue()

	if in > blk.k.CurrentVal {
		blk.greater.SetValue(true)
		blk.less.SetValue(false)
	} else if in == blk.k.CurrentVal {
		blk.greater.SetValue(false)
		blk.less.SetValue(false)
	} else {
		blk.greater.SetValue(false)
		blk.less.SetValue(true)
	}
}
