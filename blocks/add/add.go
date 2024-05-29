package add

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Add struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
}

const (
	Type  = "Add"
	Descr = "Add two inputs."

	NameIn1 = blocks.NameIn + "1"
	NameIn2 = blocks.NameIn + "2"
)

func New() blocks.Block {
	return &Add{
		in1: blocks.NewTypedInput[float64](),
		in2: blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *Add) GetType() string {
	return Type
}

func (blk *Add) GetDescription() string {
	return Descr
}

func (blk *Add) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *Add) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Add) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Add) GetWarmupPeriod() int {
	return 0
}

func (blk *Add) IsWarm() bool {
	return true
}

func (blk *Add) Init() error {
	return nil
}

func (blk *Add) Update(_ *models.Bar) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() + blk.in2.GetValue())
}
