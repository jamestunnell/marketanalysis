package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rickb777/date"
)

type Collection interface {
	GetInfo() *CollectionInfo
	GetLocation() *time.Location
	GetFirstDate() date.Date
	GetLastDate() date.Date

	IsEmpty() bool

	LoadBars(start, endIncl date.Date) (Bars, error)
	StoreBars(Bars) error
}

const (
	Resolution1Min = "1m"
)

type CollectionInfo struct {
	Symbol     string    `json:"symbol"`
	Resolution string    `json:"resolution"`
	TimeZone   string    `json:"timeZone"`
	StartDate  date.Date `json:"startDate"`
	CoreHours  *Hours    `json:"coreHours"`
}

type Hours struct {
	Open  *TimeOfDay `json:"open"`
	Close *TimeOfDay `json:"close"`
}

type TimeOfDay struct {
	Hour   int
	Minute int
}

func (tod *TimeOfDay) On(d date.Date, loc *time.Location) time.Time {
	yyyy, mm, dd := d.Date()
	return time.Date(yyyy, mm, dd, tod.Hour, tod.Minute, 0, 0, loc)
}

func (tod *TimeOfDay) MinuteOfDay() int {
	return tod.Hour*60 + tod.Minute
}

func (tod *TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", tod.Hour, tod.Minute)
}

func (tod *TimeOfDay) Parse(str string) error {
	t, err := time.Parse("15:04", str)
	if err != nil {
		return fmt.Errorf("string '%s' is not formatted as 24-hour hh:mm: %w", str, err)
	}

	tod.Hour = t.Hour()
	tod.Minute = t.Minute()

	return nil
}

func (tod *TimeOfDay) MarshalJSON() ([]byte, error) {
	return json.Marshal(tod.String())
}

func (tod *TimeOfDay) UnmarshalJSON(d []byte) error {
	var str string

	if err := json.Unmarshal(d, &str); err != nil {
		return fmt.Errorf("field is not a JSON string: %w", err)
	}

	return tod.Parse(str)
}
