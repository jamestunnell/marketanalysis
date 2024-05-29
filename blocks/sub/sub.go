package sub

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Sub struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
}

const (
	Type  = "Sub"
	Descr = "Subtract two inputs."

	NameIn1 = blocks.NameIn + "1"
	NameIn2 = blocks.NameIn + "2"
)

func New() blocks.Block {
	return &Sub{
		in1: blocks.NewTypedInput[float64](),
		in2: blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *Sub) GetType() string {
	return Type
}

func (blk *Sub) GetDescription() string {
	return Descr
}

func (blk *Sub) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *Sub) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Sub) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Sub) GetWarmupPeriod() int {
	return 0
}

func (blk *Sub) IsWarm() bool {
	return true
}

func (blk *Sub) Init() error {
	return nil
}

func (blk *Sub) Update(_ *models.Bar) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() - blk.in2.GetValue())
}
