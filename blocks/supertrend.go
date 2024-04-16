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

	atrLen *models.TypedParam[int]
	atrMul *models.TypedParam[float64]

	atr     *indicators.ATR
	prevVal *models.OHLC
}

const (
	DescrSupertrend = "Creates dynamic support and resistance levels from ATR. Uses these levels to determine trend."
	NameATRLen      = "atrLength"
	NameATRMul      = "atrMul"
	NameLower       = "lower"
	NameUpper       = "upper"
	NameTrend       = "trend"
	TypeSupertrend  = "Supertrend"
)

func NewSupertrend() models.Block {
	atrLenRange := constraints.NewValRange(2, 200)
	atrMulRange := constraints.NewValRange(0.1, 10.0)

	return &Supertrend{
		in:     models.NewTypedInput[float64](),
		trend:  models.NewTypedOutput[float64](),
		lower:  models.NewTypedOutput[float64](),
		upper:  models.NewTypedOutput[float64](),
		atrLen: models.NewParam[int](2, atrLenRange),
		atrMul: models.NewParam[float64](0.1, atrMulRange),
		atr:    nil,
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
		NameATRLen: blk.atrLen,
		NameATRMul: blk.atrMul,
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

func (blk *Supertrend) IsWarm() bool {
	return blk.atr.Warm() && (blk.prevVal != nil)
}

func (blk *Supertrend) Init() error {
	atr := indicators.NewATR(blk.atrLen.Value)

	blk.atr = atr

	return nil
}

func (blk *Supertrend) Update(bar *models.Bar) {
	defer blk.updatePrev(bar.OHLC)

	blk.atr.Update(bar.OHLC)

	if !blk.atr.Warm() || (blk.prevVal == nil) {
		return
	}

	atr := blk.atr.Current() * blk.atrMul.Value
	in := blk.in.Get()
	upper := in + atr
	lower := in - atr

	var lowerPrev float64
	if blk.lower.SetCount > 0 {
		lowerPrev = blk.lower.Value
	} else {
		lowerPrev = lower
	}

	var upperPrev float64
	if blk.upper.SetCount > 0 {
		upperPrev = blk.upper.Value
	} else {
		upperPrev = upper
	}

	if in > upperPrev {
		blk.trend.Set(1.0)
	} else if in < lowerPrev {
		blk.trend.Set(-1.0)
	} else {
		blk.trend.Set(blk.trend.Value)
	}

	if blk.prevVal.Close < upperPrev {
		blk.upper.Set(math.Min(upper, upperPrev))
	} else {
		blk.upper.Set(upper)
	}

	if blk.prevVal.Close > lowerPrev {
		blk.lower.Set(math.Max(lower, lowerPrev))
	} else {
		blk.lower.Set(lower)
	}
}

func (blk *Supertrend) updatePrev(cur *models.OHLC) {
	blk.prevVal = cur
}
