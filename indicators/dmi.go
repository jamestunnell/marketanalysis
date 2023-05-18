package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/models"
)

// DMI is a Directional Movement Index indicator.
type DMI struct {
	period                int
	prevOHLC              *models.OHLC
	mdmEMA, pdmEMA, trEMA *EMA
	pdi, mdi, dx          float64
	warm                  bool
}

func NewDMI(period int) *DMI {
	return &DMI{
		period:   period,
		prevOHLC: nil,
		mdmEMA:   NewEMA(period),
		pdmEMA:   NewEMA(period),
		trEMA:    NewEMA(period),
		pdi:      0.0,
		mdi:      0.0,
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
	pdm, mdm := pdmAndMDM(cur, dmi.prevOHLC)
	tr := TrueRange(cur, dmi.prevOHLC)

	dmi.mdmEMA.Update(mdm)
	dmi.pdmEMA.Update(pdm)
	dmi.trEMA.Update(tr)

	if !dmi.mdmEMA.Warm() {
		return
	}

	dmi.updateOutputs()

	dmi.warm = true
	dmi.prevOHLC = cur
}

func (dmi *DMI) updateOutputs() {
	dmi.pdi = dmi.pdmEMA.Current() / dmi.trEMA.Current()
	dmi.mdi = dmi.mdmEMA.Current() / dmi.trEMA.Current()
	dmi.dx = math.Abs(dmi.pdi-dmi.mdi) / (dmi.pdi + dmi.mdi)
}

// PDI returns the positive directional indicator value.
func (dmi *DMI) PDI() float64 {
	return dmi.pdi
}

// MDI returns the minus directional indicator value.
func (dmi *DMI) MDI() float64 {
	return dmi.mdi
}

// DX returns the directional movement value.
func (dmi *DMI) DX() float64 {
	return dmi.dx
}

func pdmAndMDM(cur, prev *models.OHLC) (float64, float64) {
	pdm := cur.High - prev.High
	mdm := prev.Low - cur.Low

	if pdm > mdm {
		mdm = 0.0
	} else if mdm > pdm {
		pdm = 0.0
	}

	if pdm < 0.0 {
		pdm = 0.0
	}

	if mdm < 0.0 {
		mdm = 0.0
	}

	return pdm, mdm
}
