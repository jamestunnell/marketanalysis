package predictors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util"
)

type MAOrdering struct {
	direction   models.Direction
	numPeriods  *models.TypedParam[int]
	periodStart *models.TypedParam[int]
	periodSpan  *models.TypedParam[int]
	signalLen   *models.TypedParam[int]
	threshold   *models.TypedParam[float64]
	maOrdering  *indicators.MAOrdering
}

const (
	NumPeriodsName       = "numPeriods"
	NumPeriodsMin        = 3
	NumPeriodsMaxDefault = 20

	PeriodStartName       = "periodStart"
	PeriodStartMin        = 1
	PeriodStartMaxDefault = 80

	PeriodSpanName       = "periodSpan"
	PeriodSpanMin        = 3
	PeriodSpanMaxDefault = 120

	ThresholdName = "threshold"

	TypeMAOrdering = "MAOrdering"
)

var (
	numPeriodsMax  = NewParamLimit(NumPeriodsMaxDefault)
	periodStartMax = NewParamLimit(PeriodStartMaxDefault)
	periodSpanMax  = NewParamLimit(PeriodSpanMaxDefault)
)

func init() {
	upperLimits[NumPeriodsName] = numPeriodsMax
	upperLimits[PeriodStartName] = periodStartMax
	upperLimits[PeriodSpanName] = periodSpanMax
}

func NewMAOrdering() models.Predictor {
	numPeriodsRange := constraints.NewValRange(NumPeriodsMin, numPeriodsMax.Value)
	periodStartRange := constraints.NewValRange(PeriodStartMin, periodStartMax.Value)
	periodSpanRange := constraints.NewValRange(PeriodSpanMin, periodSpanMax.Value)
	signalLenRange := constraints.NewValRange(SignalLenMin, signalLenMax.Value)
	threshRange := constraints.NewValRange(0.0, 1.0)

	return &MAOrdering{
		direction:   models.DirNone,
		numPeriods:  models.NewParam[int](numPeriodsRange),
		periodStart: models.NewParam[int](periodStartRange),
		periodSpan:  models.NewParam[int](periodSpanRange),
		signalLen:   models.NewParam[int](signalLenRange),
		threshold:   models.NewParam[float64](threshRange),
		maOrdering:  nil,
	}
}

func (mao *MAOrdering) Initialize() error {
	start := mao.periodStart.Value
	span := mao.periodSpan.Value
	periods := util.LinSpaceInts(start, start+span, mao.numPeriods.Value)

	maOrdering, err := indicators.NewMAOrdering(periods, mao.signalLen.Value)
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
		NumPeriodsName:  mao.numPeriods,
		PeriodStartName: mao.periodStart,
		PeriodSpanName:  mao.periodSpan,
		SignalLenName:   mao.signalLen,
		ThresholdName:   mao.threshold,
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

	sig := mao.maOrdering.Signal()
	thresh := mao.threshold.Value

	if sig >= thresh {
		mao.direction = models.DirUp
	} else if sig <= -thresh {
		mao.direction = models.DirDown
	} else {
		mao.direction = models.DirNone
	}
}

func (mao *MAOrdering) Direction() models.Direction {
	return mao.direction
}
