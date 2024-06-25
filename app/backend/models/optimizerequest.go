package models

import (
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/optimization"
)

type OptimizeGraphParamsRequest struct {
	JobID            string                 `json:"jobID"`
	Graph            *graph.Config          `json:"graph"`
	Symbol           string                 `json:"symbol"`
	Days             int                    `json:"days"`
	SourceQuantity   *graph.SourceQuantity  `json:"sourceQuantity"`
	TargetParams     []*graph.TargetParam   `json:"targetParams"`
	ObjectiveType    string                 `json:"objectiveType"`
	OptimizeSettings *optimization.Settings `json:"settings"`
}
