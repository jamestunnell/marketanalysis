package predictors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type MACross struct {
	direction        models.Direction
	fastPeriod       *models.TypedParam[int]
	slowPeriod       *models.TypedParam[int]
	signalLen        *models.TypedParam[int]
	fastEMA, slowEMA *indicators.EMA
	signalSMA        *indicators.SMA
}

const (
	ParamFastPeriod = "fastPeriod"
	ParamSlowPeriod = "slowPeriod"
	TypeMACross     = "MACross"

	MinPeriod = 1
)

func NewMACross() models.Predictor {
	minPeriod := constraints.NewMin(MinPeriod)

	return &MACross{
		direction:  models.DirNone,
		fastPeriod: models.NewParam[int](minPeriod),
		slowPeriod: models.NewParam[int](minPeriod),
		signalLen:  models.NewParam[int](minPeriod),
		fastEMA:    nil,
		slowEMA:    nil,
		signalSMA:  nil,
	}
}

func (mac *MACross) Type() string {
	return TypeMACross
}

func (mac *MACross) Params() models.Params {
	return models.Params{
		ParamFastPeriod: mac.fastPeriod,
		ParamSlowPeriod: mac.slowPeriod,
		ParamSignalLen:  mac.signalLen,
	}
}

func (mac *MACross) Initialize() error {
	fastPeriod := mac.fastPeriod.Value
	slowPeriod := mac.slowPeriod.Value

	if fastPeriod >= slowPeriod {
		return fmt.Errorf("fast period %d is not less than slow period %d", fastPeriod, slowPeriod)
	}

	mac.fastEMA = indicators.NewEMA(fastPeriod)
	mac.slowEMA = indicators.NewEMA(slowPeriod)
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
