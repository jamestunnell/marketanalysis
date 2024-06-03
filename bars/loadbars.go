package bars

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
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
	startDate := date.NewAt(ts.Start())
	endDate := date.NewAt(ts.End())
	primaryBars := models.Bars{}

	for d := startDate.Add(0); !d.After(endDate); d = d.Add(1) {
		dayBars, err := loader.Load(ctx, d)
		if err != nil {
			return models.Bars{}, fmt.Errorf("failed to load day bars for %s: %w", d, err)
		}

		primaryBars = append(primaryBars, dayBars.Bars...)
	}

	// can't run on this day or unknown symbol
	if len(primaryBars) == 0 {
		return models.Bars{}, nil
	}

	warmupBars := models.Bars{}

	// use bars before start time for warmup
	if startIdx, startFound := primaryBars.IndexForward(ts.Start()); startFound {
		log.Debug().Int("count", startIdx).Msg("found warmup bars before start")

		warmupBars = slices.Clone(primaryBars[:startIdx])
		primaryBars = primaryBars[startIdx:]
	}

	// ignore bars after stop time
	if endIdx, endFound := primaryBars.IndexReverse(ts.End()); endFound {
		primaryBars = primaryBars[:endIdx+1]
	}

	for warmupDate := startDate.Add(-1); len(warmupBars) < warmupPeriod; warmupDate = warmupDate.Add(-1) {
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

	runBars := append(warmupBars.LastN(warmupPeriod), primaryBars...)

	return runBars, nil
}
