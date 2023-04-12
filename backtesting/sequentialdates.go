package backtesting

import (
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type SequentialDates struct {
	current, end date.Date
}

func NewSequentialDates(dateRange timespan.DateRange) *SequentialDates {
	return &SequentialDates{
		current: dateRange.Start(),
		end:     dateRange.End(),
	}
}

func (sd *SequentialDates) Advance() {
	sd.current = sd.current.Add(1)
}

func (sd *SequentialDates) Current() date.Date {
	return sd.current
}

func (sd *SequentialDates) AnyLeft() bool {
	return sd.current.Before(sd.end)
}
