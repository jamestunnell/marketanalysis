package bars

import (
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type LoadBarsFunc func(
	symbol string,
	ts timespan.TimeSpan,
	loc *time.Location,
) (models.Bars, error)

func LoadRunBars(
	symbol string,
	ts timespan.TimeSpan,
	loc *time.Location,
	loadBars LoadBarsFunc,
	warmupPeriod int,
) (models.Bars, error) {
	primaryBars, err := loadBars(symbol, ts, loc)
	if err != nil {
		return models.Bars{}, fmt.Errorf("failed to load primary run bars: %w", err)
	}

	// can't run on this day or unknown symbol
	if len(primaryBars) == 0 {
		return models.Bars{}, nil
	}

	warmupBars := models.Bars{}
	for len(warmupBars) < warmupPeriod {
		// Move timespan backward at least one day
		ts = timespan.NewTimeSpan(
			date.NewAt(ts.Start()).Add(-1).In(loc),
			ts.Start(),
		)

		moreBars, err := loadBars(symbol, ts, loc)
		if err != nil {
			return models.Bars{}, fmt.Errorf("failed to load more bars for warmup: %w", err)
		}

		warmupBars = append(warmupBars, moreBars...)
	}

	runBars := append(sliceutils.LastN(warmupBars, warmupPeriod), primaryBars...)

	return runBars, nil
}
