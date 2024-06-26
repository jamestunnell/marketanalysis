package gate

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type Gate struct {
	in   *blocks.TypedInput[float64]
	gate *blocks.TypedInput[bool]
	out  *blocks.TypedOutput[float64]
}

const (
	Type  = "Gate"
	Descr = "Gate controlled output: input passes through when gate=true, output set to 0 when gate=false"

	NameGate = "gate"
)

func New() blocks.Block {
	return &Gate{
		in:   blocks.NewTypedInput[float64](),
		gate: blocks.NewTypedInput[bool](),
		out:  blocks.NewTypedOutput[float64](),
	}
}

func (blk *Gate) GetType() string {
	return Type
}

func (blk *Gate) GetDescription() string {
	return Descr
}

func (blk *Gate) GetParams() models.Params {
	return models.Params{}
}

func (blk *Gate) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
		NameGate:      blk.gate,
	}
}

func (blk *Gate) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Gate) GetWarmupPeriod() int {
	return 0
}

func (blk *Gate) IsWarm() bool {
	return true
}

func (blk *Gate) Init() error {
	return nil
}

func (blk *Gate) Update(cur *models.Bar, isLast bool) {
	if !blk.in.IsValueSet() || !blk.gate.IsValueSet() {
		return
	}

	if blk.gate.GetValue() {
		blk.out.SetValue(blk.in.GetValue())
	} else {
		blk.out.SetValue(0.0)
	}
}
