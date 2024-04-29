package ema

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMA struct {
	period *blocks.IntRange
	ema    *indicators.EMA
	in     *blocks.TypedInput[float64]
	out    *blocks.TypedOutput[float64]
}

const (
	Type  = "EMA"
	Descr = "Exponential Moving Average"
)

func New() blocks.Block {
	return &EMA{
		period: &blocks.IntRange{Default: 10, Min: 1, Max: 1000},
		ema:    nil,
		in:     blocks.NewTypedInput[float64](),
		out:    blocks.NewTypedOutput[float64](),
	}
}

func (blk *EMA) GetType() string {
	return Type
}

func (blk *EMA) GetDescription() string {
	return Descr
}

func (blk *EMA) GetParams() blocks.Params {
	return blocks.Params{
		blocks.NamePeriod: blk.period,
	}
}

func (blk *EMA) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *EMA) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *EMA) GetWarmupPeriod() int {
	return blk.ema.Period()
}

func (blk *EMA) IsWarm() bool {
	return blk.ema.Warm()
}

func (blk *EMA) Init() error {
	blk.ema = indicators.NewEMA(blk.period.Value)

	return nil
}

func (blk *EMA) Update(cur *models.Bar) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.ema.Update(blk.in.GetValue())

	if !blk.ema.Warm() {
		return
	}

	blk.out.SetValue(blk.ema.Current())
}