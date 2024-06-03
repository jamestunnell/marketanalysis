package bars

import (
	"context"
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
	ctx context.Context,
	symbol string,
	ts timespan.TimeSpan,
	loc *time.Location,
	loader models.DayBarsLoader,
	warmupPeriod int,
) (models.Bars, error) {
	primaryBars := models.Bars{}
	startDate := date.NewAt(ts.Start())
	endDate := date.NewAt(ts.End())

	for d := startDate; !d.After(endDate); d = d.Add(1) {
		dayBars, err := loader.Load(ctx, d)
		if err != nil {
			return models.Bars{}, fmt.Errorf("failed to load primary bars from %s: %w", d, err)
		}

		primaryBars = append(primaryBars, dayBars.Bars...)
	}

	// can't run on this day or unknown symbol
	if len(primaryBars) == 0 {
		return models.Bars{}, nil
	}

	warmupDate := startDate
	warmupBars := models.Bars{}
	for len(warmupBars) < warmupPeriod {
		// Move timespan backward at least one day
		ts = timespan.NewTimeSpan(
			date.NewAt(ts.Start()).Add(-1).In(loc),
			ts.Start(),
		)

		dayBars, err := loader.Load(ctx, warmupDate)
		if err != nil {
			return models.Bars{}, fmt.Errorf("failed to load warmup bars from %s: %w", warmupDate, err)
		}

		warmupBars = append(warmupBars, dayBars.Bars...)
	}

	runBars := append(sliceutils.LastN(warmupBars, warmupPeriod), primaryBars...)

	return runBars, nil
}
