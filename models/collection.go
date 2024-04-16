package models

import "github.com/rickb777/date/timespan"

type Collection interface {
	GetInfo() *CollectionInfo
	GetTimeSpan() timespan.TimeSpan

	LoadBars(timespan.TimeSpan) (Bars, error)
	StoreBars(Bars) error
}

const (
	Resolution1Min = "1m"
)

type CollectionInfo struct {
	Symbol     string `json:"symbol"`
	Resolution string `json:"resolution"`
}
