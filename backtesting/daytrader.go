package backtesting

import (
	"fmt"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
)

type DayTrader struct {
	Collection collection.Collection
	Strategy   models.Strategy
	Timespan   timespan.TimeSpan
	Current    date.Date
}

const (
	DayTradeMarketOpenLocalMin  = 390
	DayTradeMarketCloseLocalMin = 780
)

func NewDayTrader(c collection.Collection, s models.Strategy) Backtester {
	ts := c.Timespan()
	current := date.NewAt(ts.Start())

	return &DayTrader{
		Collection: c,
		Strategy:   s,
		Timespan:   ts,
		Current:    current,
	}
}

func (t *DayTrader) RunTest() (*Report, error) {
	open := t.Current.Local().Add(DayTradeMarketOpenLocalMin * time.Minute)
	close := t.Current.Local().Add(DayTradeMarketCloseLocalMin * time.Minute)
	ts := timespan.NewTimeSpan(open, close)
	bars := t.Collection.GetBars(ts)

	// log.Debug().
	// 	Str("current date", t.Current.FormatISO(4)).
	// 	Time("open", open).
	// 	Time("close", close).
	// 	Int("bar count", len(bars)).
	// 	Msg("day trade backtest")

	report := &Report{
		Start:     open,
		Positions: []models.Position{},
	}

	if len(bars) == 0 {
		return report, nil
	}

	err := Backtest(t.Strategy, bars)
	if err != nil {
		return nil, fmt.Errorf("failed to run backtest: %w", err)
	}

	report.Positions = t.Strategy.ClosedPositions()

	return report, nil
}

func (t *DayTrader) Advance() {
	t.Current = t.Current.Add(1)
}

func (t *DayTrader) AnyLeft() bool {
	return t.Current.UTC().Before(t.Timespan.End())
}
