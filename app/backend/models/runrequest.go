package models

import (
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/rickb777/date"
)

const (
	RunSingleDay       = "singleDay"
	RunMultiDay        = "multiDay"
	RunMultiDaySummary = "multiDaySummary"
)

type RunRequest struct {
	RunType   string        `json:"runType"`
	Graph     *graph.Config `json:"graph"`
	Symbol    string        `json:"symbol"`
	Date      date.Date     `json:"date"`
	NumCharts int           `json:"numCharts"`
}
