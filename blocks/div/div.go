package div

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Div struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]
}

const (
	Type  = "Div"
	Descr = "Divide one input by another."

	NameIn1 = blocks.NameIn + "1"
	NameIn2 = blocks.NameIn + "2"
)

func New() blocks.Block {
	return &Div{
		in1: blocks.NewTypedInput[float64](),
		in2: blocks.NewTypedInput[float64](),
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *Div) GetType() string {
	return Type
}

func (blk *Div) GetDescription() string {
	return Descr
}

func (blk *Div) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *Div) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameIn1: blk.in1,
		NameIn2: blk.in2,
	}
}

func (blk *Div) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Div) GetWarmupPeriod() int {
	return 0
}

func (blk *Div) IsWarm() bool {
	return true
}

func (blk *Div) Init() error {
	return nil
}

func (blk *Div) Update(_ *models.Bar) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	// In case the divisor is 0, Go will produce +Inf, -Inf, or NaN.
	blk.out.SetValue(blk.in1.GetValue() / blk.in2.GetValue())
}
