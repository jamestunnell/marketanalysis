package predictors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Pivot struct {
	direction models.Direction
	length    *models.TypedParam[int]
	nPivots   *models.TypedParam[int]
	pivots    *indicators.Pivots
}

const (
	NumPivotsName       = "numPivots"
	NumPivotsMaxDefault = 10
	NumPivotsMin        = 2

	PivotLenName       = "pivotLen"
	PivotLenMaxDefault = 100
	PivotLenMin        = 2

	TypePivot = "Pivot"
)

var (
	numPivotsMax = NewParamLimit(NumPivotsMaxDefault)
	pivotLenMax  = NewParamLimit(PivotLenMaxDefault)
)

func init() {
	upperLimits[PivotLenName] = pivotLenMax
	upperLimits[NumPivotsName] = numPivotsMax
}

func NewPivot() models.Predictor {
	numPivotsRange := constraints.NewValRange(NumPivotsMin, numPivotsMax.Value)
	pivotLenRange := constraints.NewValRange(PivotLenMin, pivotLenMax.Value)

	return &Pivot{
		length:  models.NewParam[int](pivotLenRange),
		nPivots: models.NewParam[int](numPivotsRange),
	}
}

func (p *Pivot) Initialize() error {
	pivots, err := indicators.NewPivots(p.length.Value, p.nPivots.Value)
	if err != nil {
		return fmt.Errorf("failed to make pivots indicator: %w", err)
	}

	p.pivots = pivots

	return nil
}

func (piv *Pivot) Type() string {
	return TypePivot
}

func (piv *Pivot) Params() models.Params {
	return models.Params{
		PivotLenName:  piv.length,
		NumPivotsName: piv.nPivots,
	}
}

func (piv *Pivot) WarmupPeriod() int {
	return piv.pivots.WarmupPeriod()
}

func (piv *Pivot) WarmUp(bars models.Bars) error {
	times := bars.Timestamps()
	closePrices := bars.ClosePrices()

	if err := piv.pivots.WarmUp(times, closePrices); err != nil {
		return fmt.Errorf("failed to warm up pivots: %w", err)
	}

	return nil
}

func (piv *Pivot) Update(bar *models.Bar) {
	if change := piv.pivots.Update(bar.Timestamp, bar.Close); change {
		piv.direction = piv.pivots.Direction()
	}
}

func (piv *Pivot) Direction() models.Direction {
	return piv.direction
}
