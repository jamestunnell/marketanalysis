package models

import (
	"github.com/jamestunnell/marketanalysis/graph"
)

const EvalSlope = "slope"

type EvalSlopeRequest struct {
	*graph.EvalSlopeConfig

	Symbol string `json:"symbol"`
}
