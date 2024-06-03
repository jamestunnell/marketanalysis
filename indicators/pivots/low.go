package pivots

import (
	"cmp"
	"slices"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type StateLow struct {
	*Pivot

	parent  *Pivots
	between []models.TimeValue[float64]
}

func NewStateLow(t time.Time, val float64, parent *Pivots) statemachine.State[float64] {
	return &StateLow{
		Pivot:   NewPivotLow(t, val),
		parent:  parent,
		between: []models.TimeValue[float64]{},
	}
}

func (state *StateLow) GetName() string {
	return "low"
}

func (state *StateLow) Enter() {
	state.parent.current = state.Pivot
}

func (state *StateLow) Exit() {
	state.parent.addCompleted(state.Pivot)
}

func (state *StateLow) Run(t time.Time, val float64) statemachine.State[float64] {
	if val <= state.Value {
		state.Value = val
		state.Timestamp = t
		state.between = []models.TimeValue[float64]{}

		return nil
	}

	if val >= state.parent.GetLastCompleted().Value {
		return NewStateHigh(t, val, state.parent)
	}

	state.between = append(state.between, models.NewTimeValue(t, val))

	if len(state.between) >= state.parent.Length {
		max := slices.MaxFunc(state.between, func(a, b models.TimeValue[float64]) int {
			return cmp.Compare(a.Value, b.Value)
		})

		return NewStateHigh(max.Time, max.Value, state.parent)
	}

	return nil
}
