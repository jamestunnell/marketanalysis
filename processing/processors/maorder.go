package processors

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util"
)

type MAOrder struct {
	numPeriods  *models.TypedParam[int]
	periodStart *models.TypedParam[int]
	periodSpan  *models.TypedParam[int]
	maOrdering  *indicators.MAOrdering
}

const (
	NumPeriodsName  = "numPeriods"
	PeriodSpanName  = "periodSpan"
	PeriodStartName = "periodStart"

	TypeMAOrder = "MAOrder"
)

func NewMAOrder() *MAOrder {
	numPeriodsRange := constraints.NewValRange(3, 15)
	startRange := constraints.NewValRange(1, 100)
	spanRange := constraints.NewValRange(5, 100)

	return &MAOrder{
		numPeriods:  models.NewParam[int](numPeriodsRange),
		periodStart: models.NewParam[int](startRange),
		periodSpan:  models.NewParam[int](spanRange),
		maOrdering:  nil,
	}
}

func (mao *MAOrder) Type() string {
	return TypeMAOrder
}

func (mao *MAOrder) Params() models.Params {
	return models.Params{
		NumPeriodsName:  mao.numPeriods,
		PeriodSpanName:  mao.periodSpan,
		PeriodStartName: mao.periodStart,
	}
}

func (mao *MAOrder) Initialize() error {
	start := mao.periodStart.Value
	span := mao.periodSpan.Value
	periods := util.LinSpaceInts(start, start+span, mao.numPeriods.Value)

	maOrdering := indicators.NewMAOrdering(periods)

	mao.maOrdering = maOrdering

	return nil
}

func (mao *MAOrder) WarmupPeriod() int {
	return mao.maOrdering.WarmupPeriod()
}

func (mao *MAOrder) Output() float64 {
	return mao.maOrdering.Correlation()
}

func (mao *MAOrder) Warm() bool {
	return mao.maOrdering.Warm()
}

func (mao *MAOrder) Update(val float64) {
	mao.maOrdering.Update(val)
}
