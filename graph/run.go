package graph

import (
	"context"
	"fmt"

	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/models"
)

func Run(
	ctx context.Context,
	cfg *Configuration,
	ts timespan.TimeSpan,
	load models.LoadBarsFunc,
) (*models.TimeSeries, error) {
	if ts.IsEmpty() {
		log.Debug().Msg("timespan is empty, returning empty time series")

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
		log.Debug().
			Stringer("start", ts.Start()).
			Stringer("end", ts.End()).
			Msg("no bars loaded, returning empty time series")

		return models.NewTimeSeries(), nil
	}

	if len(bars) <= wuPeriod {
		err := fmt.Errorf("bar count %d is not more than warmup period %d", len(bars), wuPeriod)

		return nil, err
	}

	log.Debug().
		Stringer("warmupStart", bars[0].Timestamp).
		Stringer("runStart", ts.Start()).
		Stringer("runEnd", ts.End()).
		Int("warmup bars", wuPeriod).
		Int("run bars", len(bars)-wuPeriod).
		Msgf("running model")

	for i, bar := range bars {
		g.Update(bar, i == (len(bars)-1))
	}

	timeSeries := g.GetTimeSeries()

	timeSeries.SortByTime()

	log.Debug().Msg("dropping warmup records")

	timeSeries.DropRecordsBefore(ts.Start())

	return timeSeries, nil
}
