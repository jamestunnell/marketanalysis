package predictors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type MAOrdering struct {
	direction  models.Direction
	periodMin  *models.TypedParam[int]
	periodMax  *models.TypedParam[int]
	nPeriods   *models.TypedParam[int]
	signalLen  *models.TypedParam[int]
	threshold  *models.TypedParam[float64]
	maOrdering *indicators.MAOrdering
}

const (
	ParamMinPeriod  = "minPeriod"
	ParamMaxPeriod  = "maxPeriod"
	ParamNumPeriods = "numPeriods"
	ParamSignalLen  = "signalLen"
	ParamThreshold  = "threshold"

	TypeMAOrdering = "MAOrdering"
)

func NewMAOrdering() models.Predictor {
	minPeriod := constraints.NewMin(MinPeriod)
	threshRange := constraints.NewRange(0.0, 1.0)

	return &MAOrdering{
		direction:  models.DirNone,
		periodMin:  models.NewParam[int](minPeriod),
		periodMax:  models.NewParam[int](minPeriod),
		nPeriods:   models.NewParam[int](constraints.NewMin(2)),
		signalLen:  models.NewParam[int](minPeriod),
		threshold:  models.NewParam[float64](threshRange),
		maOrdering: nil,
	}
}

func (mao *MAOrdering) Initialize() error {
	maOrdering, err := indicators.NewMAOrdering(
		mao.periodMin.Value, mao.periodMax.Value, mao.nPeriods.Value, mao.signalLen.Value)
	if err != nil {
		return fmt.Errorf("failed to make MA ordering indicator: %w", err)
	}

	mao.maOrdering = maOrdering

	return nil
}

func (mao *MAOrdering) Type() string {
	return TypeMAOrdering
}

func (mao *MAOrdering) Params() models.Params {
	return models.Params{
		ParamMinPeriod:  mao.periodMin,
		ParamMaxPeriod:  mao.periodMax,
		ParamNumPeriods: mao.nPeriods,
		ParamSignalLen:  mao.signalLen,
		ParamThreshold:  mao.threshold,
	}
}

func (mao *MAOrdering) WarmupPeriod() int {
	return mao.maOrdering.WarmupPeriod()
}

func (mao *MAOrdering) WarmUp(bars models.Bars) error {
	if err := mao.maOrdering.WarmUp(bars.ClosePrices()); err != nil {
		return fmt.Errorf("failed to warm up MA ordering indicator: %w", err)
	}

	return nil
}

func (mao *MAOrdering) Update(bar *models.Bar) {
	mao.maOrdering.Update(bar.Close)

	corr := mao.maOrdering.Correlation()
	if corr > mao.threshold.Value {
		mao.direction = models.DirUp
	} else if corr < -mao.threshold.Value {
		mao.direction = models.DirDown
	} else {
		mao.direction = models.DirNone
	}
}

func (mao *MAOrdering) Direction() models.Direction {
	return mao.direction
}
