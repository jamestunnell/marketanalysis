package atr

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type ATR struct {
	period *blocks.IntRange
	atr    *indicators.ATR
	out    *blocks.TypedOutput[float64]
}

const (
	Type  = "ATR"
	Descr = "Average True Range"
)

func New() blocks.Block {
	return &ATR{
		period: &blocks.IntRange{Default: 14, Min: 1, Max: 1000},
		atr:    nil,
		out:    blocks.NewTypedOutput[float64](),
	}
}

func (blk *ATR) GetType() string {
	return Type
}

func (blk *ATR) GetDescription() string {
	return Descr
}

func (blk *ATR) GetParams() blocks.Params {
	return blocks.Params{
		blocks.NamePeriod: blk.period,
	}
}

func (blk *ATR) GetInputs() blocks.Inputs {
	return blocks.Inputs{}
}

func (blk *ATR) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *ATR) GetWarmupPeriod() int {
	return blk.atr.Period()
}

func (blk *ATR) IsWarm() bool {
	return blk.atr.Warm()
}

func (blk *ATR) Init() error {
	blk.atr = indicators.NewATR(blk.period.Value)

	return nil
}

func (blk *ATR) Update(cur *models.Bar) {
	blk.atr.Update(cur.OHLC)

	if !blk.atr.Warm() {
		return
	}

	blk.out.SetValue(blk.atr.Current())
}
