package indicators

import (
	"github.com/jamestunnell/marketanalysis/models"
)

// EMV is an Ease of Movement Value indicator.
type EMV struct {
	emv      float64
	prevOHLC *models.OHLC
	warm     bool
}

func NewEMV() *EMV {
	return &EMV{
		prevOHLC: nil,
		emv:      0.0,
		warm:     false,
	}
}

func (emv *EMV) WarmupPeriod() int {
	return 2
}

func (emv *EMV) Warm() bool {
	return emv.warm
}

func (emv *EMV) Update(cur *models.Bar) {
	if emv.prevOHLC == nil {
		emv.prevOHLC = cur.OHLC

		return
	}

	a := (cur.High + cur.Low) / 2.0
	b := (emv.prevOHLC.High + emv.prevOHLC.Low) / 2.0
	c := float64(cur.Volume) / (cur.High - cur.Low)

	emv.emv = (a - b) / c
	emv.prevOHLC = cur.OHLC
	emv.warm = true
}

// EMV returns the ease of movement value.
func (emv *EMV) EMV() float64 {
	return emv.emv
}
