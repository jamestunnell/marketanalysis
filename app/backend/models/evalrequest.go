package models

import (
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/rickb777/date"
)

const EvalSlope = "slope"

type EvalSlopeRequest struct {
	Symbol     string         `json:"symbol"`
	Date       date.Date      `json:"date"`
	TimeZone   string         `json:"timeZone"`
	ShowWarmup bool           `json:"showWarmup,omitempty"`
	Source     *graph.Address `json:"source"`
	Predictor  *graph.Address `json:"predictor"`
	Horizon    int            `json:"horizon"`
}
