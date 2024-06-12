package models

import (
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/rickb777/date"
)

type RunDayRequest struct {
	Graph     *graph.Configuration `json:"graph"`
	Symbol    string               `json:"symbol"`
	Date      date.Date            `json:"date"`
	TimeZone  string               `json:"timeZone"`
	NumCharts int                  `json:"numCharts"`
}
