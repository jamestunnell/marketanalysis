package pivots

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type StateNeutral struct {
	*Pivot

	parent *Pivots
}

func NewStateNeutral(t time.Time, val float64, parent *Pivots) *StateNeutral {
	return &StateNeutral{
		Pivot:  NewPivotNeutral(t, val),
		parent: parent,
	}
}

func (state *StateNeutral) GetName() string {
	return "neutral"
}

func (state *StateNeutral) Enter() {
	state.parent.current = state.Pivot
}

func (state *StateNeutral) Exit() {
	state.parent.addCompleted(state.Pivot)
}

func (state *StateNeutral) Run(t time.Time, val float64) statemachine.State[float64] {
	switch {
	case val > state.Value:
		return NewStateHigh(t, val, state.parent)
	case val < state.Value:
		return NewStateLow(t, val, state.parent)
	}

	state.Timestamp = t

	return nil
}
