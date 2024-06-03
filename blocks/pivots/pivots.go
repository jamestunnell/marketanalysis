package pivots

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/blocks"
	p "github.com/jamestunnell/marketanalysis/indicators/pivots"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

type Pivots struct {
	in     *blocks.TypedInput[float64]
	out    *blocks.TypedOutputAsync[float64]
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
		out:    blocks.NewTypedOutputAsync[float64](),
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
	return blk.ind.WarmupPeriod()
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

	if !blk.ind.Update(cur.Timestamp, blk.in.GetValue()) {
		return
	}

	pivot := blk.ind.GetLastCompleted()

	log.Debug().
		Float64("value", pivot.Value).
		Time("timestamp", pivot.Timestamp).
		Time("currentTime", cur.Timestamp).
		Str("type", pivot.Type.String()).
		Msg("detected pivot point")

	blk.out.SetTimeValue(pivot.Timestamp, pivot.Value)
}
