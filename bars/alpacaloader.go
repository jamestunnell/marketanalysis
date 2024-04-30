package bars

import (
	"fmt"
	"slices"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"

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
	bars, _, err := l.getDayBars(d)

	return bars, err
}

func (l *AlpacaLoader) GetRunBars(
	d date.Date,
	wuPeriod int,
) (models.Bars, error) {
	bars, startTime, err := l.getDayBars(d)
	if err != nil {
		err = fmt.Errorf("failed to load main date bars: %w", err)

		return models.Bars{}, err
	}

	startIdx, found := slices.BinarySearchFunc(bars, startTime, func(b *models.Bar, tgt time.Time) int {
		return b.Timestamp.Compare(tgt)
	})
	if !found {
		err = fmt.Errorf("start time %s not in loaded bars", startTime)

		return models.Bars{}, err
	}

	prevDate := d.Add(-1)

	// load more warmup bars if needed
	for startIdx < wuPeriod {
		moreBars, _, err := l.getDayBars(prevDate)
		if err != nil {
			err = fmt.Errorf("failed to load more warmup bars: %w", err)

			return models.Bars{}, err
		}

		startIdx += len(moreBars)

		bars = append(moreBars, bars...)
	}

	// truncate the slice if we have enough bars for warmup
	if startIdx >= wuPeriod {
		bars = bars[startIdx-wuPeriod:]
	}

	return bars, nil
}

func (l *AlpacaLoader) getDayBars(d date.Date) (models.Bars, time.Time, error) {
	const layout = "2006-01-02T15:04"

	loc, err := time.LoadLocation(l.security.TimeZone)
	if err != nil {
		return models.Bars{}, time.Time{}, fmt.Errorf("time zone '%s' is invalid: %w", l.security.TimeZone, err)
	}

	startStr := fmt.Sprintf("%sT%s", d, l.security.Open)
	endStr := fmt.Sprintf("%sT%s", d, l.security.Close)

	start, err := time.ParseInLocation(layout, startStr, loc)
	if err != nil {
		return models.Bars{}, time.Time{}, fmt.Errorf("start time '%s' is invalid: %w", startStr, err)
	}

	end, err := time.ParseInLocation(layout, endStr, loc)
	if err != nil {
		return models.Bars{}, time.Time{}, fmt.Errorf("end time '%s' is invalid: %w", startStr, err)
	}

	ts := timespan.NewTimeSpan(start, end)

	bars, err := GetAlpacaBars(ts, l.security.Symbol, l.loc)
	if err != nil {
		return models.Bars{}, time.Time{}, err
	}

	return bars, start, nil
}
