package backend

import (
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
)

const (
	exchangesCloseOffsetMinutes = 16 * 60
	exchangesOpenOffsetMinutes  = 9*60 + 30
	tzNewYork                   = "America/New_York"
)

var locNewYork *time.Location

func init() {
	loc, err := time.LoadLocation(tzNewYork)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("timeZone", tzNewYork).
			Msg("failed to load location for New York")
	}

	locNewYork = loc
}

func GetLocationNY() *time.Location {
	return locNewYork
}

func GetCoreHours(d date.Date) timespan.TimeSpan {
	switch d.Weekday() {
	case time.Saturday, time.Sunday:
		return timespan.TimeSpan{}
	}

	start := d.In(locNewYork).Add(time.Minute * exchangesOpenOffsetMinutes)
	end := d.In(locNewYork).Add(time.Minute * exchangesCloseOffsetMinutes)

	return timespan.NewTimeSpan(start, end)
}
