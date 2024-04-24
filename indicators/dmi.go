package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/models"
)

// DMI is a Directional Movement Index indicator.
type DMI struct {
	period           int
	prevOHLC         *models.OHLC
	negDirMoveEMA    *EMA
	posDirMoveEMA    *EMA
	trueRangeEMA     *EMA
	posDI, negDI, dx float64
	warm             bool
}

func NewDMI(period int) *DMI {
	return &DMI{
		period:        period,
		prevOHLC:      nil,
		negDirMoveEMA: NewEMA(period),
		posDirMoveEMA: NewEMA(period),
		trueRangeEMA:  NewEMA(period),
		posDI:         0.0,
		negDI:         0.0,
		dx:            0.0,
		warm:          false,
	}
}

func (ind *DMI) WarmupPeriod() int {
	return 1 + ind.period
}

func (ind *DMI) Warm() bool {
	return ind.warm
}

func (ind *DMI) Update(b *models.Bar) {
	defer ind.updatePrev(b.OHLC)

	if ind.prevOHLC == nil {
		return
	}

	cur := b.OHLC
	pdm, ndm := PDMAndNDM(cur, ind.prevOHLC)
	tr := TrueRange(cur, ind.prevOHLC)

	ind.negDirMoveEMA.Update(ndm)
	ind.posDirMoveEMA.Update(pdm)
	ind.trueRangeEMA.Update(tr)

	if !ind.negDirMoveEMA.Warm() {
		return
	}

	ind.updateOutputs()

	ind.warm = true
}

func (ind *DMI) updateOutputs() {
	ind.posDI = ind.posDirMoveEMA.Current() / ind.trueRangeEMA.Current()
	ind.negDI = ind.negDirMoveEMA.Current() / ind.trueRangeEMA.Current()
	ind.dx = math.Abs(ind.posDI-ind.negDI) / (ind.posDI + ind.negDI)
}

// PDI returns the positive directional index value.
func (ind *DMI) PDI() float64 {
	return ind.posDI
}

// NDI returns the negative directional index value.
func (ind *DMI) NDI() float64 {
	return ind.negDI
}

// DX returns the directional index value.
func (ind *DMI) DX() float64 {
	return ind.dx
}

func PDMAndNDM(cur, prev *models.OHLC) (float64, float64) {
	posMove := cur.High - prev.High
	negMove := prev.Low - cur.Low

	if posMove <= negMove || posMove <= 0.0 {
		posMove = 0.0
	}

	if negMove <= posMove || negMove <= 0.0 {
		negMove = 0.0
	}

	return posMove, negMove
}

func (ind *DMI) updatePrev(cur *models.OHLC) {
	ind.prevOHLC = cur
}
