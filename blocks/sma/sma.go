package sma

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type SMA struct {
	period *models.IntParam
	sma    *indicators.SMA
	in     *blocks.TypedInput[float64]
	out    *blocks.TypedOutput[float64]
}

const (
	Descr = "Simple Moving Average"
	Type  = "SMA"
)

func New() blocks.Block {
	return &SMA{
		period: models.NewIntParam(10, models.NewGreaterEq(1)),
		sma:    nil,
		in:     blocks.NewTypedInput[float64](),
		out:    blocks.NewTypedOutput[float64](),
	}
}

func (blk *SMA) GetType() string {
	return Type
}

func (blk *SMA) GetDescription() string {
	return Descr
}

func (blk *SMA) GetParams() models.Params {
	return models.Params{
		blocks.NamePeriod: blk.period,
	}
}

func (blk *SMA) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *SMA) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *SMA) GetWarmupPeriod() int {
	return blk.sma.Period()
}

func (blk *SMA) IsWarm() bool {
	return blk.sma.Warm()
}

func (blk *SMA) Init() error {
	blk.sma = indicators.NewSMA(blk.period.CurrentVal)

	return nil
}

func (blk *SMA) Update(_ *models.Bar, isLast bool) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.sma.Update(blk.in.GetValue())

	if !blk.sma.Warm() {
		return
	}

	blk.out.SetValue(blk.sma.Current())
}
