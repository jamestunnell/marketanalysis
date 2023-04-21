package predictors

import (
	"fmt"
	"sort"

	"github.com/jamestunnell/marketanalysis/commonerrs"
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
		period1:   models.NewParam[int](period1Range),
		period2:   models.NewParam[int](period2Range),
		signalLen: models.NewParam[int](signalLenRange),
		fastEMA:   nil,
		slowEMA:   nil,
		signalSMA: nil,
	}
}

func (mac *MACross) Type() string {
	return TypeMACross
}

func (mac *MACross) Params() models.Params {
	return models.Params{
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
	wp := mac.WarmupPeriod()
	if len(bars) < wp {
		return commonerrs.NewErrMinCount("warmup bars", len(bars), wp)
	}

	vals := bars[:wp].ClosePrices()
	slowMAPeriod := mac.slowEMA.Period()
	maWUVals := vals[:slowMAPeriod]

	if err := mac.fastEMA.WarmUp(maWUVals); err != nil {
		return fmt.Errorf("failed to warm up fast EMA: %w", err)
	}

	if err := mac.slowEMA.WarmUp(maWUVals); err != nil {
		return fmt.Errorf("failed to warm up slow EMA: %w", err)
	}

	signalWUVals := make([]float64, mac.signalSMA.Period())
	for i := 0; i < mac.signalSMA.Period(); i++ {
		val := vals[i+slowMAPeriod]

		mac.fastEMA.Update(val)
		mac.slowEMA.Update(val)

		signalWUVals[i] = mac.fastEMA.Current() - mac.slowEMA.Current()
	}

	if err := mac.signalSMA.WarmUp(signalWUVals); err != nil {
		return fmt.Errorf("failed to warm up signal: %w", err)
	}

	// handle any remaining warmup vals
	for i := wp; i < len(vals); i++ {
		mac.Update(bars[i])
	}

	return nil
}

func (mac *MACross) Update(bar *models.Bar) {
	mac.fastEMA.Update(bar.Close)
	mac.slowEMA.Update(bar.Close)
	diff := mac.fastEMA.Current() - mac.slowEMA.Current()

	mac.signalSMA.Update(diff)

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
