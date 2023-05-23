package dateutils

import (
	"math/rand"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type RandomDates struct {
	rng            *rand.Rand
	start, current date.Date
	rangeDays      int32
	nLeft          int
}

func NewRandomDates(
	dateRange timespan.DateRange,
	n int,
	randSrc rand.Source) *RandomDates {
	return &RandomDates{
		start:     dateRange.Start(),
		current:   dateRange.Start(),
		rangeDays: int32(dateRange.End().Sub(dateRange.Start())),
		nLeft:     n,
		rng:       rand.New(randSrc),
	}
}

func (rd *RandomDates) Advance() {
	offset := rd.rng.Int31n(rd.rangeDays)
	rd.current = rd.start.Add(date.PeriodOfDays(offset))
	rd.nLeft--
}

func (rd *RandomDates) Current() date.Date {
	return rd.current
}

func (rd *RandomDates) AnyLeft() bool {
	return rd.nLeft > 0
}
