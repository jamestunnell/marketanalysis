package provision

import (
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type DailyBarsProvider struct {
	coll  models.Collection
	dates []date.Date
	index int
}

const (
	DayTradeMarketOpenLocalMin  = 390
	DayTradeMarketCloseLocalMin = 780
)

func NewDailyBarsProvider(
	coll models.Collection,
	dates []date.Date,
) models.BarsProvider {
	return &DailyBarsProvider{
		coll:  coll,
		dates: dates,
		index: 0,
	}
}

func (p *DailyBarsProvider) AnySetsLeft() bool {
	return p.index < len(p.dates)
}

func (p *DailyBarsProvider) Advance() {
	p.index++
}

func (p *DailyBarsProvider) CurrentSet() (models.Bars, error) {
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
