package models

import (
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/rickb777/date"
)

type BacktestRequest struct {
	Symbol    string         `json:"symbol"`
	Date      date.Date      `json:"date"`
	Predictor *graph.Address `json:"predictor"`
	Threshold float64        `json:"threshold"`
}
