package blocks

import (
	"math"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Supertrend struct {
	in                  *models.TypedInput[float64]
	trend, lower, upper *models.TypedOutput[float64]

	atrPeriod *models.TypedParam[int]
	atrMul    *models.TypedParam[float64]
	atr       *indicators.ATR
	prevVal   *models.OHLC
}

const (
	DescrSupertrend = "Creates dynamic support and resistance levels from ATR. Uses these levels to determine trend."
	NameATRPeriod   = "atrPeriod"
	NameATRMul      = "atrMul"
	NameLower       = "lower"
	NameUpper       = "upper"
	NameTrend       = "trend"
	TypeSupertrend  = "Supertrend"
)

func NewSupertrend() models.Block {
	atrPeriodRange := constraints.NewValRange(2, 200)
	atrMulRange := constraints.NewValRange(0.1, 10.0)

	return &Supertrend{
		in:        models.NewTypedInput[float64](),
		trend:     models.NewTypedOutput[float64](),
		lower:     models.NewTypedOutput[float64](),
		upper:     models.NewTypedOutput[float64](),
		atrPeriod: models.NewParam[int](2, atrPeriodRange),
		atrMul:    models.NewParam[float64](0.1, atrMulRange),
		atr:       nil,
	}
}

func (blk *Supertrend) GetType() string {
	return TypeSupertrend
}

func (blk *Supertrend) GetDescription() string {
	return DescrSupertrend
}

func (blk *Supertrend) GetParams() models.Params {
	return models.Params{
		NameATRPeriod: blk.atrPeriod,
		NameATRMul:    blk.atrMul,
	}
}

func (blk *Supertrend) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: blk.in,
	}
}

func (blk *Supertrend) GetOutputs() models.Outputs {
	return models.Outputs{
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
	blk.atr = indicators.NewATR(blk.atrPeriod.Value)

	return nil
}

func (blk *Supertrend) Update(cur *models.Bar) {
	defer blk.updatePrev(cur.OHLC)

	blk.atr.Update(cur.OHLC)

	if !blk.atr.Warm() || !blk.in.IsValueSet() {
		return
	}

	atr := blk.atr.Current() * blk.atrMul.Value
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
