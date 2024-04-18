package predictors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Pivots struct {
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

	TypePivot = "Pivots"
)

var (
	numPivotsMax = NewParamLimit(NumPivotsMaxDefault)
	pivotLenMax  = NewParamLimit(PivotLenMaxDefault)
)

func init() {
	upperLimits[PivotLenName] = pivotLenMax
	upperLimits[NumPivotsName] = numPivotsMax
}

func NewPivots() models.Predictor {
	numPivotsRange := constraints.NewValRange(NumPivotsMin, numPivotsMax.Value)
	pivotLenRange := constraints.NewValRange(PivotLenMin, pivotLenMax.Value)

	return &Pivots{
		length:  models.NewParam[int](PivotLenMin, pivotLenRange),
		nPivots: models.NewParam[int](NumPivotsMin, numPivotsRange),
	}
}

func (p *Pivots) Initialize() error {
	pivots, err := indicators.NewPivots(p.length.Value, p.nPivots.Value)
	if err != nil {
		return fmt.Errorf("failed to make pivots indicator: %w", err)
	}

	p.pivots = pivots

	return nil
}

func (p *Pivots) Type() string {
	return TypePivot
}

func (p *Pivots) Params() models.Params {
	return models.Params{
		PivotLenName:  p.length,
		NumPivotsName: p.nPivots,
	}
}

func (p *Pivots) WarmupPeriod() int {
	return p.pivots.WarmupPeriod()
}

func (p *Pivots) WarmUp(bars models.Bars) error {
	closePrices := bars.ClosePrices()

	if err := p.pivots.WarmUp(closePrices); err != nil {
		return fmt.Errorf("failed to warm up pivots: %w", err)
	}

	return nil
}

func (p *Pivots) Update(bar *models.Bar) {
	if piv, _, detected := p.pivots.Update(bar.Close); detected {
		switch piv.Type {
		case indicators.PivotHigh:
			p.direction = models.DirDown
		case indicators.PivotLow:
			p.direction = models.DirUp
		case indicators.PivotNeutral:
			p.direction = models.DirNone
		}
	}
}

func (piv *Pivots) Direction() models.Direction {
	return piv.direction
}
