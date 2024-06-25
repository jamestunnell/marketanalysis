package graph

import (
	"context"
	"fmt"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/loading"
	"github.com/jamestunnell/marketanalysis/models"
)

func RunSingleDay(
	ctx context.Context,
	cfg *Config,
	day date.Date,
	load models.LoadBarsFunc,
) (*models.TimeSeries, error) {
	ts := loading.GetCoreHours(day)

	return Run(ctx, cfg, ts, load)
}

func RunMultiDay(
	ctx context.Context,
	cfg *Config,
	dateRange timespan.DateRange,
	load models.LoadBarsFunc,
) (*models.TimeSeries, error) {
	// include entire last day by going until start of the next day
	endTime := dateRange.End().Add(1).In(loading.GetLocationNY())
	startTime := dateRange.Start().In(loading.GetLocationNY())
	ts := timespan.NewTimeSpan(startTime, endTime)

	return Run(ctx, cfg, ts, load)
}

func RunMultiDaySummary(
	ctx context.Context,
	cfg *Config,
	dateRange timespan.DateRange,
	load models.LoadBarsFunc,
) (*models.TimeSeries, error) {
	summary := models.NewTimeSeries()

	// log.Debug().
	// 	Stringer("start", startDay).
	// 	Msg("running multi-day summary")

	for d := dateRange.Start(); !d.After(dateRange.End()); d = d.Add(1) {
		timeSeries, err := RunSingleDay(ctx, cfg, d, load)
		if err != nil {
			return nil, fmt.Errorf("failed to run on day %s: %w", d, err)
		}

		endTime := loading.GetCoreHours(d).End()

		for _, q := range timeSeries.Quantities {
			for mName, mVal := range q.Measurements {
				name := q.Name + ":" + mName

				mQ, found := summary.FindQuantity(name)
				if !found {
					mQ = models.NewQuantity(name)

					summary.AddQuantity(mQ)
				}

				mQ.AddRecord(models.NewTimeValue(endTime, mVal))
			}
		}
	}

	return summary, nil
}

func Run(
	ctx context.Context,
	cfg *Config,
	ts timespan.TimeSpan,
	load models.LoadBarsFunc,
) (*models.TimeSeries, error) {
	if ts.IsEmpty() {
		log.Trace().Msg("timespan is empty, returning empty time series")

		return models.NewTimeSeries(), nil
	}

	g := New(cfg)

	if err := g.Init(); err != nil {
		return nil, fmt.Errorf("failed to init graph: %w", err)
	}

	// // re-locate timespan
	// ts = timespan.NewTimeSpan(ts.Start().In(loc), ts.End().In(loc))

	wuPeriod := g.GetWarmupPeriod()
	bars, err := models.LoadRunBars(ctx, ts, load, g.GetWarmupPeriod())
	if err != nil {
		return nil, fmt.Errorf("failed to load run bars: %w", err)
	}

	if len(bars) == 0 {
		log.Trace().
			Stringer("start", ts.Start()).
			Stringer("end", ts.End()).
			Msg("no bars loaded, returning empty time series")

		return models.NewTimeSeries(), nil
	}

	if len(bars) <= wuPeriod {
		err := fmt.Errorf("bar count %d is not more than warmup period %d", len(bars), wuPeriod)

		return nil, err
	}

	log.Trace().
		Stringer("warmupStart", bars[0].Timestamp).
		Stringer("runStart", ts.Start()).
		Stringer("runEnd", ts.End()).
		Int("warmup bars", wuPeriod).
		Int("run bars", len(bars)-wuPeriod).
		Msgf("running graph")

	for i, bar := range bars {
		g.Update(bar, i == (len(bars)-1))
	}

	timeSeries := g.GetTimeSeries()

	timeSeries.SortByTime()

	log.Trace().Msg("dropping warmup records")

	timeSeries.DropRecordsBefore(ts.Start())

	return timeSeries, nil
}
