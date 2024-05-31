package pivots

import (
	"slices"
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type StateLow struct {
	*Pivot

	parent       *Pivots
	betweenTimes []time.Time
	betweenVals  []float64
}

func NewStateLow(t time.Time, val float64, parent *Pivots) statemachine.State[float64] {
	return &StateLow{
		Pivot:        NewPivotLow(t, val),
		parent:       parent,
		betweenTimes: []time.Time{},
		betweenVals:  []float64{},
	}
}

func (state *StateLow) GetName() string {
	return "low"
}

func (state *StateLow) Enter() {
	state.parent.current = state.Pivot
}

func (state *StateLow) Exit() {
	state.parent.prev = state.Pivot
}

func (state *StateLow) Run(t time.Time, val float64) statemachine.State[float64] {
	if val <= state.Value {
		state.Value = val
		state.Timestamp = t
		state.betweenTimes = []time.Time{}
		state.betweenVals = []float64{}

		return nil
	}

	if val >= state.parent.prev.Value {
		return NewStateHigh(t, val, state.parent)
	}

	if len(state.betweenVals) >= state.parent.Length {
		minVal := slices.Min(state.betweenVals)
		minTime := state.betweenTimes[slices.Index(state.betweenVals, minVal)]

		return NewStateHigh(minTime, minVal, state.parent)
	}

	state.betweenTimes = append(state.betweenTimes, t)
	state.betweenVals = append(state.betweenVals, val)

	return nil
}
