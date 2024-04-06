package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/models"
)

// DMI is a Directional Movement Index indicator.
type DMI struct {
	period                int
	prevOHLC              *models.OHLC
	negDirMoveEMA *EMA
	posDirMoveEMA *EMA
	trueRangeEMA *EMA
	posDI, negDI, dx          float64
	warm                  bool
}

func NewDMI(period int) *DMI {
	return &DMI{
		period:   period,
		prevOHLC: nil,
		negDirMoveEMA:   NewEMA(period),
		posDirMoveEMA:   NewEMA(period),
		trueRangeEMA:    NewEMA(period),
		posDI:      0.0,
		negDI:      0.0,
		dx:       0.0,
		warm:     false,
	}
}

func (dmi *DMI) WarmupPeriod() int {
	return 1 + dmi.period
}

func (dmi *DMI) Warm() bool {
	return dmi.warm
}

func (dmi *DMI) Update(b *models.Bar) {
	if dmi.prevOHLC == nil {
		dmi.prevOHLC = b.OHLC

		return
	}

	cur := b.OHLC
	pdm, ndm := PDMAndNDM(cur, dmi.prevOHLC)
	tr := TrueRange(cur, dmi.prevOHLC)

	dmi.negDirMoveEMA.Update(ndm)
	dmi.posDirMoveEMA.Update(pdm)
	dmi.trueRangeEMA.Update(tr)

	if !dmi.negDirMoveEMA.Warm() {
		return
	}

	dmi.updateOutputs()

	dmi.warm = true
	dmi.prevOHLC = cur
}

func (dmi *DMI) updateOutputs() {
	dmi.posDI = dmi.posDirMoveEMA.Current() / dmi.trueRangeEMA.Current()
	dmi.negDI = dmi.negDirMoveEMA.Current() / dmi.trueRangeEMA.Current()
	dmi.dx = math.Abs(dmi.posDI-dmi.negDI) / (dmi.posDI + dmi.negDI)
}

// PDI returns the positive directional index value.
func (dmi *DMI) PDI() float64 {
	return dmi.posDI
}

// NDI returns the negative directional index value.
func (dmi *DMI) NDI() float64 {
	return dmi.negDI
}

// DX returns the directional index value.
func (dmi *DMI) DX() float64 {
	return dmi.dx
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
