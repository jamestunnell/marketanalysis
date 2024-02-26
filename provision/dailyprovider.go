package provision

import (
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type DailyBarProvider struct {
	coll  collection.Collection
	dates []date.Date
	index int
}

const (
	DayTradeMarketOpenLocalMin  = 390
	DayTradeMarketCloseLocalMin = 780
)

func EachNonEmptyBarSet(
	provider BarProvider,
	eachSet func(bars models.Bars)) error {
	for provider.AnySetsLeft() {
		bars, err := provider.CurrentSet()
		if err != nil {
			return fmt.Errorf("failed to get current bar set: %w", err)
		}

		if len(bars) != 0 {
			eachSet(bars)
		}

		provider.Advance()
	}

	return nil
}

func NewDailyBarProvider(coll collection.Collection, dates []date.Date) BarProvider {
	return &DailyBarProvider{
		coll:  coll,
		dates: dates,
		index: 0,
	}
}

func (p *DailyBarProvider) AnySetsLeft() bool {
	return p.index < len(p.dates)
}

func (p *DailyBarProvider) Advance() {
	p.index++
}

func (p *DailyBarProvider) CurrentSet() (models.Bars, error) {
	d := p.dates[p.index]
	dayStart := d.Local()
	nextDayStart := dayStart.Add(24 * time.Hour)
	ts := timespan.NewTimeSpan(dayStart, nextDayStart)

	bars, err := p.coll.LoadBars(ts)
	if err != nil {
		return models.Bars{}, fmt.Errorf("failed to load bars: %w", err)
	}

	open := d.Local().Add(DayTradeMarketOpenLocalMin * time.Minute)
	close := d.Local().Add(DayTradeMarketCloseLocalMin * time.Minute)

	bars = sliceutils.Where(bars, func(b *models.Bar) bool {
		return b.Timestamp.After(open) && b.Timestamp.Before(close)
	})

	return bars, nil
}
