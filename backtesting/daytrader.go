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
	Collection     collection.Collection
	Strategy       models.Strategy
	DateController DateController
}

type DateController interface {
	Current() date.Date
	Advance()
	AnyLeft() bool
}

const (
	DayTradeMarketOpenLocalMin  = 390
	DayTradeMarketCloseLocalMin = 780
)

func NewDayTrader(
	c collection.Collection,
	s models.Strategy,
	dc DateController) *DayTrader {
	return &DayTrader{
		Collection:     c,
		Strategy:       s,
		DateController: dc,
	}
}

func (t *DayTrader) RunTest() (*Report, error) {
	open := t.DateController.Current().Local().Add(DayTradeMarketOpenLocalMin * time.Minute)
	close := t.DateController.Current().Local().Add(DayTradeMarketCloseLocalMin * time.Minute)
	ts := timespan.NewTimeSpan(open, close)
	bars := t.Collection.GetBars(ts)

	report := &Report{
		Start:     open,
		Positions: models.Positions{},
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
	t.DateController.Advance()
}

func (t *DayTrader) AnyLeft() bool {
	return t.DateController.AnyLeft()
}
