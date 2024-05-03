package bars

import (
	"fmt"
	"slices"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/models"
)

type AlpacaLoader struct {
	security *models.Security
	loc      *time.Location
}

func NewAlpacaLoader(s *models.Security) Loader {
	return &AlpacaLoader{
		security: s,
		loc:      nil,
	}
}

func (l *AlpacaLoader) Init() error {
	loc, err := time.LoadLocation(l.security.TimeZone)
	if err != nil {
		return fmt.Errorf("time zone '%s' is invalid: %w", l.security.TimeZone, err)
	}

	l.loc = loc

	return nil
}

func (l *AlpacaLoader) GetBars(ts timespan.TimeSpan) (models.Bars, error) {
	return GetAlpacaBars(ts, l.security.Symbol, l.loc)
}

func (l *AlpacaLoader) GetDayBars(d date.Date) (models.Bars, error) {
	bars, err := l.getDayBars(d)

	return bars, err
}

func (l *AlpacaLoader) GetRunBars(
	d date.Date,
	wuPeriod int,
) (models.Bars, error) {
	const layout = "2006-01-02T15:04"

	bars, err := l.getDayBars(d)
	if err != nil {
		err = fmt.Errorf("failed to load main date bars: %w", err)

		return models.Bars{}, err
	}

	if len(bars) == 0 {
		return models.Bars{}, nil
	}

	// runs begin at open time-of-day
	startStr := fmt.Sprintf("%sT%s", d, l.security.Open)
	start, err := time.ParseInLocation(layout, startStr, l.loc)
	if err != nil {
		err = fmt.Errorf("failed to parse start time %s: %w", startStr, err)

		return models.Bars{}, err
	}

	// runs ends at close time-of-day
	endStr := fmt.Sprintf("%sT%s", d, l.security.Close)
	end, err := time.ParseInLocation(layout, endStr, l.loc)
	if err != nil {
		err = fmt.Errorf("failed to parse end time %s: %w", endStr, err)

		return models.Bars{}, err
	}

	endIdx, found := slices.BinarySearchFunc(bars, end, compareBarByTimestamp)
	if found || endIdx > 0 {
		nTrim := endIdx - len(bars)

		bars = bars[:endIdx]

		log.Debug().Msgf("alpacaloader: trimed %d bars at end", nTrim)
	} else if endIdx == 0 {
		err = fmt.Errorf("end time %s is less than all loaded bars", end)

		return models.Bars{}, err
	}

	startIdx, found := slices.BinarySearchFunc(bars, start, compareBarByTimestamp)
	if !found && startIdx == len(bars) {
		err = fmt.Errorf("start time %s is more than all loaded bars", start)

		return models.Bars{}, err
	}

	prevDate := d.Add(-1)

	// load more warmup bars if needed
	for startIdx < wuPeriod {
		moreBars, err := l.getDayBars(prevDate)
		if err != nil {
			err = fmt.Errorf("failed to load more warmup bars: %w", err)

			return models.Bars{}, err
		}

		startIdx += len(moreBars)

		bars = append(moreBars, bars...)
		prevDate = prevDate.Add(-1)
	}

	// truncate the slice if we have enough bars for warmup
	if startIdx >= wuPeriod {
		bars = bars[startIdx-wuPeriod:]
	}

	return bars, nil
}

func compareBarByTimestamp(b *models.Bar, tgt time.Time) int {
	return b.Timestamp.Compare(tgt)
}

func (l *AlpacaLoader) getDayBars(d date.Date) (models.Bars, error) {
	loc, err := time.LoadLocation(l.security.TimeZone)
	if err != nil {
		return models.Bars{}, fmt.Errorf("time zone '%s' is invalid: %w", l.security.TimeZone, err)
	}

	start := d.In(loc)
	end := d.Add(1).In(loc).Add(-time.Minute)
	ts := timespan.NewTimeSpan(start, end)

	bars, err := GetAlpacaBars(ts, l.security.Symbol, l.loc)
	if err != nil {
		return models.Bars{}, err
	}

	return bars, nil
}
