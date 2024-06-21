package multitrend

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type MultiTrend2 struct {
	in1 *blocks.TypedInput[float64]
	in2 *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]

	thresh      *blocks.FloatParam
	votesNeeded *blocks.IntParam
}

const (
	TypeMultitrend2  = "Multitrend2"
	DescrMultitrend2 = `Aggregates two trend inputs (uptrend=1, downtrend=-1)`

	NameVotesNeeded = "votesNeeded"
	NameThresh      = "thresh"
)

func NewMultitrend2() blocks.Block {
	return &MultiTrend2{
		thresh:      blocks.NewFloatParam(0.375, blocks.NewRangeExcl(0.0, 1.0)),
		votesNeeded: blocks.NewIntParam(1, blocks.NewOneOf([]int{1, 2})),
		in1:         blocks.NewTypedInput[float64](),
		in2:         blocks.NewTypedInput[float64](),
		out:         blocks.NewTypedOutput[float64](),
	}
}

func (blk *MultiTrend2) GetType() string {
	return TypeMultitrend2
}

func (blk *MultiTrend2) GetDescription() string {
	return DescrMultitrend2
}

func (blk *MultiTrend2) GetParams() blocks.Params {
	return blocks.Params{
		NameThresh:      blk.thresh,
		NameVotesNeeded: blk.votesNeeded,
	}
}

func (blk *MultiTrend2) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn + "1": blk.in1,
		blocks.NameIn + "2": blk.in2,
	}
}

func (blk *MultiTrend2) GetOutputs() blocks.Outputs {
	return blocks.Outputs{blocks.NameOut: blk.out}
}

func (blk *MultiTrend2) GetWarmupPeriod() int {
	return 0
}

func (blk *MultiTrend2) IsWarm() bool {
	return true
}

func (blk *MultiTrend2) Init() error {
	return nil
}

func (blk *MultiTrend2) Update(cur *models.Bar, isLast bool) {
	if !blk.in1.IsValueSet() || !blk.in2.IsValueSet() {
		return
	}

	in1 := blk.in1.GetValue()
	in2 := blk.in2.GetValue()

	dir := 0

	if in1 > blk.thresh.CurrentVal {
		dir += 1
	} else if in1 < -blk.thresh.CurrentVal {
		dir -= 1
	}

	if in2 > blk.thresh.CurrentVal {
		dir += 1
	} else if in2 < -blk.thresh.CurrentVal {
		dir -= 1
	}

	if dir > 0 {
		blk.out.SetValue(1.0)
	} else if dir < 0 {
		blk.out.SetValue(-1.0)
	} else {
		blk.out.SetValue(0.0)
	}
}
