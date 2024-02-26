package provision

import (
	"math"
	"math/rand"
	"time"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rickb777/date"
)

func SplitCollectionRandomly(
	coll collection.Collection,
	split float64,
	randSource rand.Source) (trainingBars, testingBars BarProvider, err error) {
	dateRange := coll.TimeSpan().DateRangeIn(time.Local)
	nDays := int(dateRange.Days())

	numTraining := int(math.Round(split * float64(nDays)))
	if numTraining < 1 || numTraining >= nDays {
		return nil, nil, commonerrs.NewErrOutOfRange("training days", numTraining, 1, nDays-1)
	}

	days := sliceutils.New(nDays, func(idx int) int {
		return idx
	})
	rng := rand.New(randSource)

	rng.Shuffle(nDays, func(i, j int) {
		days[i], days[j] = days[j], days[i]
	})

	start := dateRange.Start()
	trainingDates := sliceutils.Map(days[:numTraining], func(day int) date.Date {
		return start.Add(date.PeriodOfDays(day))
	})
	testingDates := sliceutils.Map(days[numTraining:], func(day int) date.Date {
		return start.Add(date.PeriodOfDays(day))
	})

	trainingBars = NewDailyBarProvider(coll, trainingDates)
	testingBars = NewDailyBarProvider(coll, testingDates)

	return
}
