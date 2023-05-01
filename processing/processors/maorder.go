package processors

import (
	"sort"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/montanaflynn/stats"
	"golang.org/x/exp/slices"
)

type MAOrder struct {
	numPeriods  *models.TypedParam[int]
	periodStart *models.TypedParam[int]
	periodSpan  *models.TypedParam[int]

	MAs            []*indicators.EMA
	warmupPeriod   int
	uptrendIndices []float64
	correlation    float64
	nPeriods       int
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
		numPeriods:     models.NewParam[int](numPeriodsRange),
		periodStart:    models.NewParam[int](startRange),
		periodSpan:     models.NewParam[int](spanRange),
		MAs:            []*indicators.EMA{},
		uptrendIndices: []float64{},
		correlation:    0.0,
		nPeriods:       0,
		warmupPeriod:   0,
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
	mao.MAs = make([]*indicators.EMA, len(periods))
	mao.uptrendIndices = make([]float64, len(periods))

	for i, period := range periods {
		mao.MAs[i] = indicators.NewEMA(period)
		mao.uptrendIndices[i] = float64(i)
	}

	mao.warmupPeriod = sliceutils.Last(periods)
	mao.correlation = 0.0
	mao.nPeriods = len(periods)

	return nil
}

func (mao *MAOrder) WarmupPeriod() int {
	return mao.warmupPeriod
}

func (mao *MAOrder) Output() float64 {
	return mao.correlation
}

func (mao *MAOrder) WarmUp(vals []float64) {
	for _, ma := range mao.MAs {
		_ = ma.WarmUp(vals)
	}
}

func (mao *MAOrder) Update(val float64) {
	for _, ma := range mao.MAs {
		ma.Update(val)
	}

	mao.updateCorrelation()
}

func (mao *MAOrder) updateCorrelation() {
	values := sliceutils.Map(mao.MAs, func(ma *indicators.EMA) float64 {
		return ma.Current()
	})

	sort.Float64s(values)

	// find the indices of the values after they are sorted
	currentIndices := sliceutils.Map(mao.MAs, func(ma *indicators.EMA) float64 {
		return float64(slices.Index(values, ma.Current()))
	})

	// how well do they correlate with the perfect uptrend ordering?
	a, _ := stats.Correlation(mao.uptrendIndices, currentIndices)

	mao.correlation = a
}
