package provision

import (
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type DailyBars struct {
	Collection collection.Collection
	Date       date.Date
	bars       models.Bars
	index      int
}

func NewDailyBarSet(c collection.Collection, d date.Date) *DailyBars {
	return &DailyBars{
		Collection: c,
		Date:       d,
		bars:       models.Bars{},
		index:      0,
	}
}

func (db *DailyBars) Initialize() error {
	dayStart := db.Date.Local()
	dayEnd := dayStart.Add(time.Hour * 24)
	ts := timespan.NewTimeSpan(dayStart, dayEnd)

	bars, err := db.Collection.LoadBars(ts)
	if err != nil {
		return fmt.Errorf("failed to load bars: %w", err)
	}

	db.bars = bars
	db.index = 0

	return nil
}

func (db *DailyBars) EachBar(each func(bar *models.Bar)) {
	for _, bar := range db.bars {
		each(bar)
	}
}
