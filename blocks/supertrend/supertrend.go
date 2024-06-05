package supertrend

import (
	"math"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Supertrend struct {
	in                  *blocks.TypedInput[float64]
	trend, lower, upper *blocks.TypedOutput[float64]

	atrPeriod *blocks.TypedParam[int]
	atrMul    *blocks.TypedParam[float64]
	atr       *indicators.ATR
	prevVal   *models.OHLC
}

const (
	Type  = "Supertrend"
	Descr = "Creates dynamic support and resistance levels from ATR. Uses these levels to determine trend."

	NameATRPeriod = "atrPeriod"
	NameATRMul    = "atrMul"
	NameLower     = "lower"
	NameUpper     = "upper"
	NameTrend     = "trend"
)

func New() blocks.Block {
	return &Supertrend{
		in:        blocks.NewTypedInput[float64](),
		trend:     blocks.NewTypedOutput[float64](),
		lower:     blocks.NewTypedOutput[float64](),
		upper:     blocks.NewTypedOutput[float64](),
		atrPeriod: blocks.NewTypedParam(20, blocks.NewGreaterEqual(1)),
		atrMul:    blocks.NewTypedParam(5.0, blocks.NewGreater(0.0)),
		atr:       nil,
	}
}

func (blk *Supertrend) GetType() string {
	return Type
}

func (blk *Supertrend) GetDescription() string {
	return Descr
}

func (blk *Supertrend) GetParams() blocks.Params {
	return blocks.Params{
		NameATRPeriod: blk.atrPeriod,
		NameATRMul:    blk.atrMul,
	}
}

func (blk *Supertrend) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *Supertrend) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		NameLower: blk.lower,
		NameUpper: blk.upper,
		NameTrend: blk.trend,
	}
}

func (blk *Supertrend) GetWarmupPeriod() int {
	return blk.atr.Period()
}

func (blk *Supertrend) IsWarm() bool {
	return blk.atr.Warm()
}

func (blk *Supertrend) Init() error {
	blk.atr = indicators.NewATR(blk.atrPeriod.CurrentVal)

	return nil
}

func (blk *Supertrend) Update(cur *models.Bar) {
	defer blk.updatePrev(cur.OHLC)

	blk.atr.Update(cur.OHLC)

	if !blk.atr.Warm() || !blk.in.IsValueSet() {
		return
	}

	atr := blk.atr.Current() * blk.atrMul.CurrentVal
	in := blk.in.GetValue()
	upper := in + atr
	lower := in - atr

	lowerPrev := lower
	if blk.lower.IsValueSet() {
		lowerPrev = blk.lower.GetValue()
	}

	upperPrev := upper
	if blk.upper.IsValueSet() {
		upperPrev = blk.upper.GetValue()
	}

	if !blk.trend.IsValueSet() {
		blk.trend.SetValue(0.0)
	}

	if blk.trend.GetValue() <= 0.0 && in > upperPrev {
		blk.trend.SetValue(1.0)
	} else if blk.trend.GetValue() >= 0.0 && in < lowerPrev {
		blk.trend.SetValue(-1.0)
	}

	if blk.prevVal.Close < upperPrev {
		blk.upper.SetValue(math.Min(upper, upperPrev))
	} else {
		blk.upper.SetValue(upper)
	}

	if blk.prevVal.Close > lowerPrev {
		blk.lower.SetValue(math.Max(lower, lowerPrev))
	} else {
		blk.lower.SetValue(lower)
	}
}

func (blk *Supertrend) updatePrev(cur *models.OHLC) {
	blk.prevVal = cur
}
