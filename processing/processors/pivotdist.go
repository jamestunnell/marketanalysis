package processors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type PivotDist struct {
	output   float64
	pivots   *indicators.Pivots
	pivotLen *models.TypedParam[int]
}

const (
	PivotLenName = "pivotLen"

	TypePivotDist = "PivotDist"
)

func NewPivotDist() *PivotDist {
	lenRange := constraints.NewValRange(1, 200)

	return &PivotDist{
		output:   0.0,
		pivots:   nil,
		pivotLen: models.NewParam[int](lenRange),
	}
}

func (pd *PivotDist) Type() string {
	return TypePivotDist
}

func (pd *PivotDist) Params() models.Params {
	return models.Params{
		PivotLenName: pd.pivotLen,
	}
}

func (pd *PivotDist) Initialize() error {
	pivs, err := indicators.NewPivots(pd.pivotLen.Value, 4)
	if err != nil {
		return fmt.Errorf("failed to make pivots indicator: %w", err)
	}

	pd.pivots = pivs
	pd.output = 0.0

	return nil
}

func (pd *PivotDist) WarmupPeriod() int {
	return pd.pivots.WarmupPeriod()
}

func (pd *PivotDist) Output() float64 {
	return pd.output
}

func (pd *PivotDist) WarmUp(vals []float64) error {
	if err := pd.pivots.WarmUp(vals); err != nil {
		return fmt.Errorf("failed to warm up pivots indicator: %w", err)
	}

	return nil
}

func (pd *PivotDist) Update(val float64) {
	// TODO
}
