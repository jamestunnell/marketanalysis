package models

import (
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type Collection interface {
	GetInfo() *CollectionInfo
	IsEmpty() bool
	GetLastDate() date.Date

	LoadBars(timespan.TimeSpan) (Bars, error)
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
}
