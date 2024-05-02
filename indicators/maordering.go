package indicators

import (
	"sort"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/montanaflynn/stats"
	"golang.org/x/exp/slices"
)

type MAOrdering struct {
	MAs            []*EMA
	orderCurrent   []float64
	orderUptrend   []float64
	valuesSorted   []float64
	valuesUnsorted []float64
	correlation    float64
	nPeriods       int
	lastMA         *EMA
}

const MinNumPeriods = 2

func NewMAOrdering(periods []int) *MAOrdering {
	slices.Sort(periods)

	nPeriods := len(periods)
	mas := make([]*EMA, nPeriods)
	orderUptrend := make([]float64, nPeriods)
	lastIndex := nPeriods - 1

	for i, period := range periods {
		mas[i] = NewEMA(period)
		orderUptrend[lastIndex-i] = float64(i) // uptrend means lowest period MA has the highest value
	}

	ind := &MAOrdering{
		MAs:            mas,
		orderUptrend:   orderUptrend,
		orderCurrent:   make([]float64, nPeriods),
		valuesSorted:   make([]float64, nPeriods),
		valuesUnsorted: make([]float64, nPeriods),
		correlation:    0.0,
		nPeriods:       nPeriods,
		lastMA:         sliceutils.Last(mas),
	}

	return ind
}

func (ind *MAOrdering) WarmupPeriod() int {
	return ind.lastMA.Period()
}

func (ind *MAOrdering) Warm() bool {
	return ind.lastMA.Warm()
}

func (ind *MAOrdering) Update(val float64) {
	for _, ma := range ind.MAs {
		ma.Update(val)
	}

	if ind.Warm() {
		ind.updateCorrelation()
	}
}

func (ind *MAOrdering) Correlation() float64 {
	return ind.correlation
}

func (ind *MAOrdering) updateCorrelation() {
	for i, ma := range ind.MAs {
		ind.valuesUnsorted[i] = ma.Current()
		ind.valuesSorted[i] = ma.Current()
	}

	sort.Float64s(ind.valuesSorted)

	for i, val := range ind.valuesSorted {
		ind.orderCurrent[i] = float64(slices.Index(ind.valuesUnsorted, val))
	}

	// how well do they correlate with the perfect uptrend ordering?
	a, _ := stats.Correlation(ind.orderUptrend, ind.orderCurrent)

	ind.correlation = a
}
