package bar

import (
	"github.com/rickb777/date/timespan"
)

func BarsTimespan(bars []*Bar) timespan.TimeSpan {
	if len(bars) == 0 {
		return timespan.TimeSpan{}
	}

	min := bars[0].Timestamp
	max := bars[0].Timestamp

	for i := 1; i < len(bars); i++ {
		if bars[i].Timestamp.Before(min) {
			min = bars[i].Timestamp
		}

		if bars[i].Timestamp.After(max) {
			max = bars[i].Timestamp
		}
	}

	return timespan.NewTimeSpan(min, max)
}
