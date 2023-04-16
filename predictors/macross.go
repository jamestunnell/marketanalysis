package predictors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type MACross struct {
	direction        models.Direction
	params           models.Params
	fastEMA, slowEMA *indicators.EMA
	signal           *indicators.SMA
}

const (
	ParamFastPeriod = "fastPeriod"
	ParamSlowPeriod = "slowPeriod"
	TypeMACross     = "MACross"
)

func NewMACross(params models.Params) (models.Predictor, error) {
	fastPeriod, err := params.GetInt(ParamFastPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get fast period param: %w", err)
	}

	slowPeriod, err := params.GetInt(ParamSlowPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get slow period param: %w", err)
	}

	signalLen, err := params.GetInt(ParamSignalLen)
	if err != nil {
		return nil, fmt.Errorf("failed to get signal len param: %w", err)
	}

	if fastPeriod >= slowPeriod {
		return nil, fmt.Errorf("fast period %d is not less than slow period %d", fastPeriod, slowPeriod)
	}

	fastEMA := indicators.NewEMA(fastPeriod)
	slowEMA := indicators.NewEMA(slowPeriod)
	signal := indicators.NewSMA(signalLen)

	mac := &MACross{
		direction: models.DirNone,
		params:    params,
		fastEMA:   fastEMA,
		slowEMA:   slowEMA,
		signal:    signal,
	}

	return mac, nil
}

func (mac *MACross) Type() string {
	return TypeMACross
}

func (mac *MACross) Params() models.Params {
	return mac.params
}

func (mac *MACross) WarmupPeriod() int {
	return mac.slowEMA.Period() + mac.signal.Period()
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

	signalWUVals := make([]float64, mac.signal.Period())
	for i := 0; i < mac.signal.Period(); i++ {
		val := vals[i+slowMAPeriod]

		mac.fastEMA.Update(val)
		mac.slowEMA.Update(val)

		signalWUVals[i] = mac.fastEMA.Current() - mac.slowEMA.Current()
	}

	if err := mac.signal.WarmUp(signalWUVals); err != nil {
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

	mac.signal.Update(diff)

	sig := mac.signal.Current()

	if sig > 0.0 {
		mac.direction = models.DirUp
	} else if sig < 0.0 {
		mac.direction = models.DirDown
	}
}

func (mac *MACross) Direction() models.Direction {
	return mac.direction
}
