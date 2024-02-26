package provision

import (
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type DailyBarSeq struct {
	Collection collection.Collection
	Date       date.Date
}

type DailyBarSeqs struct {
	Collection collection.Collection
	Dates      []date.Date
}

func NewDailyBarSeq(c collection.Collection, d date.Date) *DailyBarSeq {
	return &DailyBarSeq{
		Collection: c,
		Date:       d,
	}
}

func NewDailyBarSeqs(c collection.Collection, dates ...date.Date) *DailyBarSeqs {
	return &DailyBarSeqs{
		Collection: c,
		Dates:      dates,
	}
}

func (db *DailyBarSeq) EachBar(each func(bar *models.Bar) error) error {
	dayStart := db.Date.Local()
	dayEnd := dayStart.Add(time.Hour * 24)
	ts := timespan.NewTimeSpan(dayStart, dayEnd)

	bars, err := db.Collection.LoadBars(ts)
	if err != nil {
		return fmt.Errorf("failed to load bars: %w", err)
	}

	for _, bar := range bars {
		if err = each(bar); err != nil {
			return fmt.Errorf("bar handler failed: %w", err)
		}
	}

	return nil
}

func (db *DailyBarSeqs) EachSequence(each func(seq BarSequence) error) error {
	for _, d := range db.Dates {
		seq := NewDailyBarSeq(db.Collection, d)

		if err := each(seq); err != nil {
			return fmt.Errorf("seq handler failed: %w", err)
		}
	}

	return nil
}
