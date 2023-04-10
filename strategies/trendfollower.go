package strategies

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type TrendFollower struct {
	direction              int
	params                 models.Params
	fastPeriod, slowPeriod int
	fastEMA, slowEMA       models.Indicator
	closedPositions        []models.Position
	openPosition           models.Position
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

	fastEMA := indicators.NewEMA(fastPeriod)
	slowEMA := indicators.NewEMA(slowPeriod)

	tf := &TrendFollower{
		direction:       0,
		closedPositions: []models.Position{},
		openPosition:    nil,
		params:          params,
		fastEMA:         fastEMA,
		slowEMA:         slowEMA,
		fastPeriod:      fastPeriod,
		slowPeriod:      slowPeriod,
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
	return tf.slowPeriod
}

func (tf *TrendFollower) WarmUp(bars []*models.Bar) error {
	fastWUBars := bars[len(bars)-tf.fastPeriod:]
	if err := tf.fastEMA.WarmUp(fastWUBars); err != nil {
		return fmt.Errorf("failed to warm up fast EMA: %w", err)
	}

	if err := tf.slowEMA.WarmUp(bars); err != nil {
		return fmt.Errorf("failed to warm up slow EMA: %w", err)
	}

	return nil
}

func (tf *TrendFollower) Update(bar *models.Bar) {
	fast := tf.fastEMA.Update(bar)
	slow := tf.slowEMA.Update(bar)
	diff := fast - slow

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
