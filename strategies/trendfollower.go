package strategies

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type TrendFollower struct {
	direction        int
	params           models.Params
	fastEMA, slowEMA *indicators.EMA
	closedPositions  []models.Position
	openPosition     models.Position
}

const (
	ParamFastPeriod   = "fastPeriod"
	ParamSlowPeriod   = "slowPeriod"
	TypeTrendFollower = "TrendFollower"
)

func NewTrendFollower(params models.Params) (models.Strategy, error) {
	fastPeriod, err := params.GetInt(ParamFastPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get fast period param: %w", err)
	}

	slowPeriod, err := params.GetInt(ParamSlowPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get slow period param: %w", err)
	}

	if fastPeriod >= slowPeriod {
		return nil, fmt.Errorf("fast period %d is not less than slow period %d", fastPeriod, slowPeriod)
	}

	fastEMA := indicators.NewEMA(fastPeriod)
	slowEMA := indicators.NewEMA(slowPeriod)

	tf := &TrendFollower{
		direction:       0,
		closedPositions: []models.Position{},
		openPosition:    nil,
		params:          params,
		fastEMA:         fastEMA,
		slowEMA:         slowEMA,
	}

	return tf, nil
}

func (tf *TrendFollower) Type() string {
	return TypeTrendFollower
}

func (tf *TrendFollower) Params() models.Params {
	return tf.params
}

func (tf *TrendFollower) ClosedPositions() []models.Position {
	return tf.closedPositions
}

func (tf *TrendFollower) Close(bar *models.Bar) {
	if tf.openPosition != nil {
		tf.openPosition.Close(bar.Timestamp, bar.Close)

		tf.closedPositions = append(tf.closedPositions, tf.openPosition)
		tf.openPosition = nil
	}
}

func (tf *TrendFollower) WarmupPeriod() int {
	return tf.slowEMA.Period()
}

func (tf *TrendFollower) Initialize(bars models.Bars) error {
	fastWUBars := bars[:tf.fastEMA.Period()]
	if err := tf.fastEMA.WarmUp(fastWUBars.ClosePrices()); err != nil {
		return fmt.Errorf("failed to warm up fast EMA: %w", err)
	}

	for i := tf.fastEMA.Period(); i < len(bars); i++ {
		tf.fastEMA.Update(bars[i].Close)
	}

	if err := tf.slowEMA.WarmUp(bars.ClosePrices()); err != nil {
		return fmt.Errorf("failed to warm up slow EMA: %w", err)
	}

	return nil
}

func (tf *TrendFollower) Update(bar *models.Bar) {
	tf.fastEMA.Update(bar.Close)
	tf.slowEMA.Update(bar.Close)
	diff := tf.fastEMA.Current() - tf.slowEMA.Current()

	if tf.openPosition == nil {
		if diff > 0.0 {
			tf.openPosition = models.NewLongPosition(bar.Timestamp, bar.Close)
			tf.direction = 1
		} else {
			tf.openPosition = models.NewShortPosition(bar.Timestamp, bar.Close)
			tf.direction = -1
		}

		return
	}

	if diff > 0.0 && tf.direction == -1 {
		tf.Close(bar)

		tf.openPosition = models.NewLongPosition(bar.Timestamp, bar.Close)
		tf.direction = 1
	} else if diff < 0.0 && tf.direction == 1 {
		tf.Close(bar)

		tf.openPosition = models.NewShortPosition(bar.Timestamp, bar.Close)
		tf.direction = -1
	}
}
