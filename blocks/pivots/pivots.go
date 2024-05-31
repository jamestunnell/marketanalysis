package pivots

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/blocks"
	p "github.com/jamestunnell/marketanalysis/indicators/pivots"
	"github.com/jamestunnell/marketanalysis/models"
)

type Pivots struct {
	in     *blocks.TypedInput[float64]
	out    *blocks.TypedOutput[float64]
	length *blocks.IntRange
	ind    *p.Pivots
}

const (
	Type  = "Pivots"
	Descr = "Pivot points ."
)

func New() blocks.Block {
	return &Pivots{
		in:     blocks.NewTypedInput[float64](),
		out:    blocks.NewTypedOutput[float64](),
		length: &blocks.IntRange{Default: 15, Min: 2, Max: 1000},
		ind:    nil,
	}
}

func (blk *Pivots) GetType() string {
	return Type
}

func (blk *Pivots) GetDescription() string {
	return Descr
}

func (blk *Pivots) GetParams() blocks.Params {
	return blocks.Params{
		blocks.NameLength: blk.length,
	}
}

func (blk *Pivots) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *Pivots) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *Pivots) GetWarmupPeriod() int {
	return 0
}

func (blk *Pivots) IsWarm() bool {
	return true
}

func (blk *Pivots) Init() error {
	ind, err := p.New(blk.length.Value)
	if err != nil {
		return fmt.Errorf("failed to make pivots indicator: %w", err)
	}

	blk.ind = ind

	return nil
}

func (blk *Pivots) Update(cur *models.Bar) {
	if !blk.in.IsValueSet() {
		return
	}

	tvIn := blk.in.GetTimeValue()

	if !blk.ind.Update(tvIn.Time, tvIn.Value) {
		return
	}

	pivot := blk.ind.GetLatest()

	blk.out.SetTimeValue(pivot.Timestamp, pivot.Value)
}
