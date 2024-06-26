package aroon

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Aroon struct {
	period *blocks.IntParam
	aroon  *indicators.Aroon
	in     *blocks.TypedInput[float64]
	diff   *blocks.TypedOutput[float64]
	dn     *blocks.TypedOutput[float64]
	up     *blocks.TypedOutput[float64]
}

const (
	Type  = "Aroon"
	Descr = "Aroon indicator identifies trend changes in the price of an asset, as well as the strength of that trend."

	NameDiff = "diff"
	NameDown = "down"
	NameUp   = "up"
)

func New() blocks.Block {
	return &Aroon{
		period: blocks.NewIntParam(10, blocks.NewGreaterEqual(1)),
		aroon:  nil,
		in:     blocks.NewTypedInput[float64](),
		up:     blocks.NewTypedOutput[float64](),
		dn:     blocks.NewTypedOutput[float64](),
		diff:   blocks.NewTypedOutput[float64](),
	}
}

func (blk *Aroon) GetType() string {
	return Type
}

func (blk *Aroon) GetDescription() string {
	return Descr
}

func (blk *Aroon) GetParams() blocks.Params {
	return blocks.Params{
		blocks.NamePeriod: blk.period,
	}
}

func (blk *Aroon) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *Aroon) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		NameDiff: blk.diff,
		NameDown: blk.dn,
		NameUp:   blk.up,
	}
}

func (blk *Aroon) GetWarmupPeriod() int {
	return blk.aroon.WarmupPeriod()
}

func (blk *Aroon) IsWarm() bool {
	return blk.aroon.Warm()
}

func (blk *Aroon) Init() error {
	blk.aroon = indicators.NewAroon(blk.period.CurrentVal)

	return nil
}

func (blk *Aroon) Update(_ *models.Bar, isLast bool) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.aroon.Update(blk.in.GetValue())

	if !blk.aroon.Warm() {
		return
	}

	blk.diff.SetValue(blk.aroon.Up() - blk.aroon.Down())
	blk.dn.SetValue(blk.aroon.Down())
	blk.up.SetValue(blk.aroon.Up())
}
