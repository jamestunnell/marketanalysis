package indicators

import (
	"fmt"
	"sort"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/montanaflynn/stats"
	"golang.org/x/exp/slices"
)

type MAOrdering struct {
	// params         models.Params
	MAs            []*EMA
	signal         *SMA
	warmupPeriod   int
	uptrendIndices []float64
	correlation    float64
	nPeriods       int
}

const MinNumPeriods = 2

func NewMAOrdering(
	periodMin, periodMax, nPeriods, signalLen int) (*MAOrdering, error) {
	if periodMin >= periodMax {
		return nil, fmt.Errorf("period min %d is not less than period max %d", periodMin, periodMax)
	}

	if nPeriods < MinNumPeriods {
		return nil, fmt.Errorf("n periods %d is less than %d", nPeriods, MinNumPeriods)
	}

	periods := util.LinSpaceInts(periodMin, periodMax, nPeriods)
	mas := make([]*EMA, nPeriods)
	uptrendIndices := make([]float64, nPeriods)

	for i, period := range periods {
		mas[i] = NewEMA(period)
		uptrendIndices[i] = float64(i)
	}

	signalSMA := NewSMA(signalLen)
	wuPeriod := signalSMA.Period() + sliceutils.Last(mas).Period()

	mao := &MAOrdering{
		// params:         params,
		signal:         signalSMA,
		MAs:            mas,
		warmupPeriod:   wuPeriod,
		uptrendIndices: uptrendIndices,
		correlation:    0.0,
		nPeriods:       nPeriods,
	}

	return mao, nil
}

func (mao *MAOrdering) WarmupPeriod() int {
	return mao.warmupPeriod
}

func (mao *MAOrdering) WarmUp(vals []float64) error {
	if len(vals) < mao.warmupPeriod {
		return commonerrs.NewErrMinCount("warmup values", len(vals), mao.warmupPeriod)
	}

	// warm up the MAs
	lastMAPeriod := sliceutils.Last(mao.MAs).Period()
	maWUVals := vals[:lastMAPeriod]
	for _, ma := range mao.MAs {
		err := ma.WarmUp(maWUVals)
		if err != nil {
			return fmt.Errorf("failed to warm up MA(%d): %w", ma.Period(), err)
		}
	}

	signalWUVals := make([]float64, mao.signal.Period())
	for i := 0; i < mao.signal.Period(); i++ {
		val := vals[i+lastMAPeriod]
		for _, ma := range mao.MAs {
			ma.Update(val)
		}

		mao.updateCorrelation()

		signalWUVals[i] = mao.correlation
	}

	if err := mao.signal.WarmUp(signalWUVals); err != nil {
		return fmt.Errorf("failed to warm up signal: %w", err)
	}

	// handle any remaining warmup vals
	for i := mao.warmupPeriod; i < len(vals); i++ {
		mao.Update(vals[i])
	}

	return nil
}

func (mao *MAOrdering) Update(val float64) {
	for _, ma := range mao.MAs {
		ma.Update(val)
	}

	mao.updateCorrelation()

	mao.signal.Update(val)
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
