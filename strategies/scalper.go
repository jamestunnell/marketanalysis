package strategies

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Scalper struct {
	direction              int
	params                 models.Params
	fastPeriod, slowPeriod int
	takeProfit             float64
	fastEMA, slowEMA       *indicators.EMA
	closedPositions        []models.Position
	openPosition           models.Position
}

const (
	TypeScalper     = "Scalper"
	ParamTakeProfit = "takeProfit"
)

func NewScalper(params models.Params) (models.Strategy, error) {
	fastPeriod, err := params.GetInt(ParamFastPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get fast period param: %w", err)
	}

	slowPeriod, err := params.GetInt(ParamSlowPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get slow period param: %w", err)
	}

	takeProfit, err := params.GetFloat(ParamTakeProfit)
	if err != nil {
		return nil, fmt.Errorf("failed to get take profit param: %w", err)
	}

	if fastPeriod >= slowPeriod {
		return nil, fmt.Errorf("fast period %d is not less than slow period %d", fastPeriod, slowPeriod)
	}

	fastEMA := indicators.NewEMA(fastPeriod)
	slowEMA := indicators.NewEMA(slowPeriod)

	sc := &Scalper{
		direction:       0,
		closedPositions: []models.Position{},
		openPosition:    nil,
		params:          params,
		fastEMA:         fastEMA,
		slowEMA:         slowEMA,
		fastPeriod:      fastPeriod,
		slowPeriod:      slowPeriod,
		takeProfit:      takeProfit,
	}

	return sc, nil
}

func (sc *Scalper) Type() string {
	return TypeScalper
}

func (sc *Scalper) Params() models.Params {
	return sc.params
}

func (sc *Scalper) ClosedPositions() []models.Position {
	return sc.closedPositions
}

func (sc *Scalper) Close(bar *models.Bar) {
	if sc.openPosition != nil {
		sc.openPosition.Close(bar.Timestamp, bar.Close)

		sc.closedPositions = append(sc.closedPositions, sc.openPosition)
		sc.openPosition = nil
	}
}

func (sc *Scalper) WarmupPeriod() int {
	return sc.slowPeriod
}

func (sc *Scalper) Initialize(bars models.Bars) error {
	fastWUBars := bars[:sc.fastPeriod]
	if err := sc.fastEMA.WarmUp(fastWUBars.ClosePrices()); err != nil {
		return fmt.Errorf("failed to warm up fast EMA: %w", err)
	}

	for i := sc.fastPeriod; i < len(bars); i++ {
		sc.fastEMA.Update(bars[i].Close)
	}

	if err := sc.slowEMA.WarmUp(bars.ClosePrices()); err != nil {
		return fmt.Errorf("failed to warm up slow EMA: %w", err)
	}

	sc.closedPositions = []models.Position{}
	sc.openPosition = nil

	return nil
}

func (sc *Scalper) Update(bar *models.Bar) {
	sc.fastEMA.Update(bar.Close)
	sc.slowEMA.Update(bar.Close)

	diff := sc.fastEMA.Current() - sc.slowEMA.Current()

	if sc.openPosition == nil {
		if diff > 0.0 && sc.direction <= 0 {
			sc.openPosition = models.NewLongPosition(bar.Timestamp, bar.Close)
			sc.direction = 1
		} else if diff < 0.0 && sc.direction >= 0 {
			sc.openPosition = models.NewShortPosition(bar.Timestamp, bar.Close)
			sc.direction = -1
		}

		return
	}

	pl, _ := sc.openPosition.OpenProfitLoss(bar.Close)
	if pl >= sc.takeProfit {
		sc.Close(bar)

		sc.openPosition = nil
	} else if diff > 0.0 && sc.direction == -1 {
		sc.Close(bar)

		sc.openPosition = models.NewLongPosition(bar.Timestamp, bar.Close)
		sc.direction = 1
	} else if diff < 0.0 && sc.direction == 1 {
		sc.Close(bar)

		sc.openPosition = models.NewShortPosition(bar.Timestamp, bar.Close)
		sc.direction = -1
	}
}
