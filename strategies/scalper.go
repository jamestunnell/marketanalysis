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
	fastEMA, slowEMA       models.Indicator
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

	tf := &Scalper{
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

	return tf, nil
}

func (tf *Scalper) Type() string {
	return TypeScalper
}

func (tf *Scalper) Params() models.Params {
	return tf.params
}

func (tf *Scalper) ClosedPositions() []models.Position {
	return tf.closedPositions
}

func (tf *Scalper) Close(bar *models.Bar) {
	if tf.openPosition != nil {
		tf.openPosition.Close(bar.Timestamp, bar.Close)

		tf.closedPositions = append(tf.closedPositions, tf.openPosition)
		tf.openPosition = nil
	}
}

func (tf *Scalper) WarmupPeriod() int {
	return tf.slowPeriod
}

func (tf *Scalper) WarmUp(bars []*models.Bar) error {
	fastWUBars := bars[len(bars)-tf.fastPeriod:]
	if err := tf.fastEMA.WarmUp(fastWUBars); err != nil {
		return fmt.Errorf("failed to warm up fast EMA: %w", err)
	}

	if err := tf.slowEMA.WarmUp(bars); err != nil {
		return fmt.Errorf("failed to warm up slow EMA: %w", err)
	}

	return nil
}

func (tf *Scalper) Update(bar *models.Bar) {
	fast := tf.fastEMA.Update(bar)
	slow := tf.slowEMA.Update(bar)
	diff := fast - slow

	if tf.openPosition == nil {
		if diff > 0.0 && tf.direction <= 0 {
			tf.openPosition = models.NewLongPosition(bar.Timestamp, bar.Close)
			tf.direction = 1
		} else if diff < 0.0 && tf.direction >= 0 {
			tf.openPosition = models.NewShortPosition(bar.Timestamp, bar.Close)
			tf.direction = -1
		}

		return
	}

	pl, _ := tf.openPosition.OpenProfitLoss(bar.Close)
	if pl >= tf.takeProfit {
		tf.Close(bar)

		tf.openPosition = nil
	} else if diff > 0.0 && tf.direction == -1 {
		tf.Close(bar)

		tf.openPosition = models.NewLongPosition(bar.Timestamp, bar.Close)
		tf.direction = 1
	} else if diff < 0.0 && tf.direction == 1 {
		tf.Close(bar)

		tf.openPosition = models.NewShortPosition(bar.Timestamp, bar.Close)
		tf.direction = -1
	}
}
