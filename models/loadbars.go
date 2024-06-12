package models

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
)

type LoadBarsFunc func(ctx context.Context, d date.Date) (Bars, error)

func LoadRunBars(
	ctx context.Context,
	symbol string,
	ts timespan.TimeSpan,
	loc *time.Location,
	load LoadBarsFunc,
	warmupPeriod int,
) (Bars, error) {
	startDate := date.NewAt(ts.Start())
	endDate := date.NewAt(ts.End())
	primaryBars := Bars{}

	for d := startDate.Add(0); !d.After(endDate); d = d.Add(1) {
		bars, err := load(ctx, d)
		if err != nil {
			return Bars{}, fmt.Errorf("failed to load day bars for %s: %w", d, err)
		}

		primaryBars = append(primaryBars, bars...)
	}

	// can't run on this day or unknown symbol
	if len(primaryBars) == 0 {
		return Bars{}, nil
	}

	warmupBars := Bars{}

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

		bars, err := load(ctx, warmupDate)
		if err != nil {
			return Bars{}, fmt.Errorf("failed to load warmup bars from %s: %w", warmupDate, err)
		}

		warmupBars = append(warmupBars, bars...)
	}

	runBars := append(warmupBars.LastN(warmupPeriod), primaryBars...)

	return runBars, nil
}
