package predictors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Pivot struct {
	direction models.Direction
	params    models.Params
	pivots    *indicators.Pivots
}

const (
	ParamPeriod  = "length"
	ParamNPivots = "nPivots"
	TypePivot    = "Pivot"
)

func NewPivot(params models.Params) (models.Predictor, error) {
	length, err := params.GetInt(ParamPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get length param: %w", err)
	}

	nPivots, err := params.GetInt(ParamNPivots)
	if err != nil {
		return nil, fmt.Errorf("failed to get nPivots param: %w", err)
	}

	pivots, err := indicators.NewPivots(length, nPivots)
	if err != nil {
		return nil, fmt.Errorf("failed to make pivots indicator: %w", err)
	}

	piv := &Pivot{
		direction: models.DirNone,
		params:    params,
		pivots:    pivots,
	}

	return piv, nil
}

func (piv *Pivot) Type() string {
	return TypePivot
}

func (piv *Pivot) Params() models.Params {
	return piv.params
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
