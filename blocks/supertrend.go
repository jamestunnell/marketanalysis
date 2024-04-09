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
	NameATRLen     = "atrLength"
	NameATRMul     = "atrMul"
	NameLower      = "lower"
	NameUpper      = "upper"
	NameTrend      = "trend"
	TypeSupertrend = "Supertrend"
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

func (sup *Supertrend) GetType() string {
	return TypeSupertrend
}

func (sup *Supertrend) GetDescription() string {
	return DescrSupertrend
}

func (sup *Supertrend) GetParams() models.Params {
	return models.Params{
		NameATRLen: sup.atrLen,
		NameATRMul: sup.atrMul,
	}
}

func (sup *Supertrend) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: sup.in,
	}
}

func (sup *Supertrend) GetOutputs() models.Outputs {
	return models.Outputs{
		NameLower: sup.lower,
		NameUpper: sup.upper,
		NameTrend: sup.trend,
	}
}

func (sup *Supertrend) IsWarm() bool {
	return sup.atr.Warm() && (sup.prevVal != nil)
}

func (sup *Supertrend) Init() error {
	atr := indicators.NewATR(sup.atrLen.Value)

	sup.atr = atr

	return nil
}

func (sup *Supertrend) Update(bar *models.Bar) {
	sup.atr.Update(bar.OHLC)

	if !sup.atr.Warm() || (sup.prevVal == nil) {
		sup.prevVal = bar.OHLC

		return
	}

	atr := sup.atr.Current() * sup.atrMul.Value
	up := sup.in.Get() + atr
	dn := sup.in.Get() - atr

	if bar.Close > sup.upper.Value {
		sup.trend.Set(1.0)
	} else if bar.Close < sup.lower.Value {
		sup.trend.Set(-1.0)
	}

	if sup.prevVal.Close < sup.upper.Value {
		sup.upper.Set(math.Min(up, sup.upper.Value))
	} else {
		sup.upper.Set(up)
	}

	if sup.prevVal.Close > sup.lower.Value {
		sup.upper.Set(math.Max(dn, sup.lower.Value))
	} else {
		sup.lower.Set(dn)
	}

	sup.prevVal = bar.OHLC
}
