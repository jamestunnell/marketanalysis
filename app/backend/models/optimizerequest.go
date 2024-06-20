package models

import (
	"github.com/jamestunnell/marketanalysis/graph"
)

type OptimizeRequest struct {
	Graph            *graph.Config           `json:"graph"`
	Symbol           string                  `json:"symbol"`
	Days             int                     `json:"days"`
	SourceQuantity   *graph.SourceQuantity   `json:"sourceQuantity"`
	TargetParams     []*graph.TargetParam    `json:"targetParams"`
	OptimizeSettings *graph.OptimizeSettings `json:"settings"`
}
