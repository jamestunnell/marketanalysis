package mul

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Mul struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
}

const (
	Type  = "Mul"
	Descr = "Multiply two inputs."

	NameIn1 = blocks.NameIn + "1"
	NameIn2 = blocks.NameIn + "2"
)

func New() blocks.Block {
	return &Mul{
		in1: blocks.NewTypedInput[float64](),
		in2: blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *Mul) GetType() string {
	return Type
}

func (blk *Mul) GetDescription() string {
	return Descr
}

func (blk *Mul) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *Mul) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Mul) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Mul) GetWarmupPeriod() int {
	return 0
}

func (blk *Mul) IsWarm() bool {
	return true
}

func (blk *Mul) Init() error {
	return nil
}

func (blk *Mul) Update(_ *models.Bar) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	blk.out.SetValue(blk.in1.GetValue() * blk.in2.GetValue())
}
