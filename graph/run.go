package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/recorders"
)

const (
	exchangesCloseOffsetMinutes = 16 * 60
	exchangesOpenOffsetMinutes  = 9*60 + 30
	exchangesTZ                 = "America/New_York"
)

var exchangeLoc *time.Location

func init() {
	loc, err := time.LoadLocation(exchangesTZ)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("timeZone", exchangesTZ).
			Msg("failed to load exchange location")
	}

	exchangeLoc = loc
}

func RunDay(
	ctx context.Context,
	cfg *Configuration,
	symbol string,
	d date.Date,
	loc *time.Location,
	loader models.DayBarsLoader,
) (*models.TimeSeries, error) {
	return Run(ctx, cfg, symbol, GetCoreHours(d), loc, loader)
}

func Run(
	ctx context.Context,
	cfg *Configuration,
	symbol string,
	ts timespan.TimeSpan,
	loc *time.Location,
	loader models.DayBarsLoader,
) (*models.TimeSeries, error) {
	if ts.IsEmpty() {
		log.Debug().Msg("timespan is empty, returning empty time series")

		return models.NewTimeSeries(), nil
	}

	g := New(cfg)
	r := recorders.NewTimeSeries(loc)

	if err := g.Init(r); err != nil {
		return nil, fmt.Errorf("failed to init graph: %w", err)
	}

	wuPeriod := g.GetWarmupPeriod()
	bars, err := bars.LoadRunBars(ctx, symbol, ts, loc, loader, g.GetWarmupPeriod())
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
		Stringer("start", ts.Start()).
		Stringer("end", ts.End()).
		Int("warmup bars", wuPeriod).
		Int("run bars", len(bars)-wuPeriod).
		Msgf("running model")

	for _, bar := range bars {
		g.Update(bar)
	}

	if err = r.Finalize(); err != nil {
		return nil, fmt.Errorf("failed to finalize recording: %w", err)
	}

	r.DropRecordsBefore(ts.Start())

	return r.TimeSeries, nil
}

func GetCoreHours(d date.Date) timespan.TimeSpan {
	const layout = "2006-01-02T15:04"

	switch d.Weekday() {
	case time.Saturday, time.Sunday:
		return timespan.TimeSpan{}
	}

	start := d.In(exchangeLoc).Add(time.Minute * exchangesOpenOffsetMinutes)
	end := d.In(exchangeLoc).Add(time.Minute * exchangesCloseOffsetMinutes)

	return timespan.NewTimeSpan(start, end)
}
