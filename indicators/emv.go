package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

// EMV is an Ease of Movement Value indicator.
type EMV struct {
	emv      float64
	prevOHLC *models.OHLC
}

func NewEMV() *EMV {
	return &EMV{
		prevOHLC: nil,
		emv:      0.0,
	}
}

func (emv *EMV) WarmupPeriod() int {
	return 2
}

func (emv *EMV) WarmUp(bars models.Bars) error {
	n := len(bars)
	if n < 2 {
		return commonerrs.NewErrMinCount("warmup bars", n, 2)
	}

	emv.emv = 0.0
	emv.prevOHLC = bars[0].OHLC

	for i := 1; i < n; i++ {
		emv.Update(bars[i])
	}

	return nil
}

func (emv *EMV) Update(cur *models.Bar) {
	a := (cur.High + cur.Low) / 2.0
	b := (emv.prevOHLC.High + emv.prevOHLC.Low) / 2.0
	c := float64(cur.Volume) / (cur.High - cur.Low)

	emv.emv = (a - b) / c
	emv.prevOHLC = cur.OHLC
}

// EMV returns the ease of movement value.
func (emv *EMV) EMV() float64 {
	return emv.emv
}
