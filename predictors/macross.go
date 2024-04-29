package predictors

import (
	"sort"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type MACross struct {
	direction        models.Direction
	period1          *models.TypedParam[int]
	period2          *models.TypedParam[int]
	signalLen        *models.TypedParam[int]
	fastEMA, slowEMA *indicators.EMA
	signalSMA        *indicators.SMA
}

const (
	Period1Name       = "period1"
	Period1MaxDefault = 200
	Period1Min        = 1

	Period2Name       = "period2"
	Period2MaxDefault = 200
	Period2Min        = 1

	SignalLenName       = "signalLen"
	SignalLenMaxDefault = 20
	SignalLenMin        = 1

	TypeMACross = "MACross"
)

var (
	period1Max   = NewParamLimit(Period1MaxDefault)
	period2Max   = NewParamLimit(Period2MaxDefault)
	signalLenMax = NewParamLimit(SignalLenMaxDefault)
)

func init() {
	upperLimits[Period1Name] = period1Max
	upperLimits[Period2Name] = period2Max
	upperLimits[SignalLenName] = signalLenMax
}

func NewMACross() models.Predictor {
	period1Range := constraints.NewValRange(Period1Min, period1Max.Value)
	period2Range := constraints.NewValRange(Period2Min, period2Max.Value)
	signalLenRange := constraints.NewValRange(SignalLenMin, period2Max.Value)

	return &MACross{
		direction: models.DirNone,
		period1:   models.NewParam[int](Period1Min, period1Range),
		period2:   models.NewParam[int](Period2Min, period2Range),
		signalLen: models.NewParam[int](SignalLenMin, signalLenRange),
		fastEMA:   nil,
		slowEMA:   nil,
		signalSMA: nil,
	}
}

func (mac *MACross) Type() string {
	return TypeMACross
}

func (mac *MACross) Params() blocks.Params {
	return blocks.Params{
		Period1Name:   mac.period1,
		Period2Name:   mac.period2,
		SignalLenName: mac.signalLen,
	}
}

func (mac *MACross) Initialize() error {
	periods := []int{
		mac.period1.Value,
		mac.period2.Value}

	sort.Ints(periods)

	mac.fastEMA = indicators.NewEMA(periods[0])
	mac.slowEMA = indicators.NewEMA(periods[1])
	mac.signalSMA = indicators.NewSMA(mac.signalLen.Value)

	return nil
}

func (mac *MACross) WarmupPeriod() int {
	return mac.slowEMA.Period() + mac.signalSMA.Period()
}

func (mac *MACross) WarmUp(bars models.Bars) error {
	for _, bar := range bars {
		mac.Update(bar)
	}

	return nil
}

func (mac *MACross) Update(bar *models.Bar) {
	mac.fastEMA.Update(bar.Close)
	mac.slowEMA.Update(bar.Close)

	if !mac.slowEMA.Warm() {
		return
	}

	diff := mac.fastEMA.Current() - mac.slowEMA.Current()

	mac.signalSMA.Update(diff)

	if !mac.signalSMA.Warm() {
		return
	}

	sig := mac.signalSMA.Current()

	if sig > 0.0 {
		mac.direction = models.DirUp
	} else if sig < 0.0 {
		mac.direction = models.DirDown
	}
}

func (mac *MACross) Direction() models.Direction {
	return mac.direction
}
