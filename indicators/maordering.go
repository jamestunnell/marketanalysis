package indicators

import (
	"sort"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/montanaflynn/stats"
	"golang.org/x/exp/slices"
)

type MAOrdering struct {
	MAs            []*EMA
	uptrendIndices []float64
	correlation    float64
	nPeriods       int
	lastMA         *EMA
}

const MinNumPeriods = 2

func NewMAOrdering(periods []int) *MAOrdering {
	mas := make([]*EMA, len(periods))
	uptrendIndices := make([]float64, len(periods))

	for i, period := range periods {
		mas[i] = NewEMA(period)
		uptrendIndices[i] = float64(i)
	}

	mao := &MAOrdering{
		MAs:            mas,
		uptrendIndices: uptrendIndices,
		correlation:    0.0,
		nPeriods:       len(periods),
		lastMA:         sliceutils.Last(mas),
	}

	return mao
}

func (mao *MAOrdering) WarmupPeriod() int {
	return mao.lastMA.Period()
}

func (mao *MAOrdering) Warm() bool {
	return mao.lastMA.Warm()
}

func (mao *MAOrdering) Update(val float64) {
	for _, ma := range mao.MAs {
		ma.Update(val)
	}

	if mao.Warm() {
		mao.updateCorrelation()
	}
}

func (mao *MAOrdering) Correlation() float64 {
	return mao.correlation
}

func (mao *MAOrdering) updateCorrelation() {
	values := sliceutils.Map(mao.MAs, func(ma *EMA) float64 {
		return ma.Current()
	})

	sort.Float64s(values)

	// find the indices of the values after they are sorted
	currentIndices := sliceutils.Map(mao.MAs, func(ma *EMA) float64 {
		return float64(slices.Index(values, ma.Current()))
	})

	// how well do they correlate with the perfect uptrend ordering?
	a, _ := stats.Correlation(mao.uptrendIndices, currentIndices)

	mao.correlation = a
}
