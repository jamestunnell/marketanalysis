package pivots

import (
	"slices"
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type StateHigh struct {
	*Pivot

	parent       *Pivots
	betweenTimes []time.Time
	betweenVals  []float64
}

func NewStateHigh(t time.Time, val float64, parent *Pivots) statemachine.State[float64] {
	return &StateHigh{
		Pivot:        NewPivotHigh(t, val),
		parent:       parent,
		betweenTimes: []time.Time{},
		betweenVals:  []float64{},
	}
}

func (state *StateHigh) GetName() string {
	return "high"
}

func (state *StateHigh) Enter() {
	state.parent.current = state.Pivot
}

func (state *StateHigh) Exit() {
	state.parent.prev = state.Pivot
}

func (state *StateHigh) Run(t time.Time, val float64) statemachine.State[float64] {
	if val >= state.Pivot.Value {
		state.Pivot.Value = val
		state.Pivot.Timestamp = t
		state.betweenTimes = []time.Time{}
		state.betweenVals = []float64{}

		return nil
	}

	if val <= state.parent.prev.Value {
		return NewStateLow(t, val, state.parent)
	}

	if len(state.betweenVals) >= state.parent.Length {
		minVal := slices.Min(state.betweenVals)
		minTime := state.betweenTimes[slices.Index(state.betweenVals, minVal)]

		return NewStateLow(minTime, minVal, state.parent)
	}

	state.betweenTimes = append(state.betweenTimes, t)
	state.betweenVals = append(state.betweenVals, val)

	return nil
}
