package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

// EMV is an Ease of Movement Value indicator.
type EMV struct {
	emv      float64
	prevOHLC *models.OHLC
	scale    float64
	sma      *SMA
}

func NewEMV(period int, scale float64) (*EMV, error) {
	if scale <= 0.0 {
		return nil, commonerrs.NewErrNotPositive("scale", scale)
	}

	emv := &EMV{
		prevOHLC: nil,
		emv:      0.0,
		scale:    scale,
		sma:      NewSMA(period),
	}

	return emv, nil
}

func (ind *EMV) WarmupPeriod() int {
	return 1 + ind.sma.Period()
}

func (ind *EMV) PartlyWarm() bool {
	return ind.prevOHLC != nil
}

func (ind *EMV) FullyWarm() bool {
	return ind.sma.Warm()
}

func (ind *EMV) Update(cur *models.Bar) {
	defer ind.udpatePrev(cur)

	if ind.prevOHLC == nil {
		return
	}

	distMoved := ((cur.High + cur.Low) / 2.0) - ((ind.prevOHLC.High + ind.prevOHLC.Low) / 2.0)
	boxRatio := (float64(cur.Volume) / ind.scale) / (cur.High - cur.Low)

	ind.emv = distMoved / boxRatio

	ind.sma.Update(ind.emv)
}

// Curent returns the current 1-period ease of movement value.
func (ind *EMV) Current() float64 {
	return ind.emv
}

// Average returns the averaged EMV.
func (ind *EMV) Average() float64 {
	return ind.sma.Current()
}

func (ind *EMV) udpatePrev(cur *models.Bar) {
	ind.prevOHLC = cur.OHLC
}
