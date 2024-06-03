package pivots

import (
	"cmp"
	"slices"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type StateHigh struct {
	*Pivot

	parent  *Pivots
	between []models.TimeValue[float64]
}

func NewStateHigh(t time.Time, val float64, parent *Pivots) statemachine.State[float64] {
	return &StateHigh{
		Pivot:   NewPivotHigh(t, val),
		parent:  parent,
		between: []models.TimeValue[float64]{},
	}
}

func (state *StateHigh) GetName() string {
	return "high"
}

func (state *StateHigh) Enter() {
	state.parent.current = state.Pivot
}

func (state *StateHigh) Exit() {
	state.parent.addCompleted(state.Pivot)
}

func (state *StateHigh) Run(t time.Time, val float64) statemachine.State[float64] {
	if val >= state.Value {
		state.Value = val
		state.Timestamp = t
		state.between = []models.TimeValue[float64]{}

		return nil
	}

	if val <= state.parent.GetLastCompleted().Value {
		return NewStateLow(t, val, state.parent)
	}

	state.between = append(state.between, models.NewTimeValue(t, val))

	if len(state.between) >= state.parent.Length {
		min := slices.MinFunc(state.between, func(a, b models.TimeValue[float64]) int {
			return cmp.Compare(a.Value, b.Value)
		})

		return NewStateLow(min.Time, min.Value, state.parent)
	}

	return nil
}
